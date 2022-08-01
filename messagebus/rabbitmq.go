package messagebus

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/rid"
	"github.com/inconshreveable/log15"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rotisserie/eris"
)

const (
	// DefaultQueueExpire value for queue auto expire
	DefaultQueueExpire = 30000
	// DefaultMessageExpire value message expire
	DefaultMessageExpire = 30000

	// exchange names
	exchangeEvent    = "ari.event"
	exchangePing     = "ari.ping"
	exchangeAnnounce = "ari.announce"
	exchangeRequest  = "ari.request"

	// type of identifiers
	ridConsumer    = "co"
	ridConsumerReq = "rp"
	ridQueue       = "q"
	ridCorrelation = "cr"
)

// RabbitmqBus is MessageBus implementation for RabbitMQ
type RabbitmqBus struct {
	Config Config
	Log    log15.Logger

	conn          *amqp091.Connection
	channel       *amqp091.Channel
	countTimeouts int64
	isClosed      bool
	mu            sync.RWMutex
}

// OptionRabbitmqFunc options for RabbitMQ
type OptionRabbitmqFunc func(n *RabbitmqBus)

// NewRabbitmqBus creates a RabbitmqBus
func NewRabbitmqBus(config Config, options ...OptionRabbitmqFunc) *RabbitmqBus {

	mbus := RabbitmqBus{
		Config: config,
	}

	for _, optfn := range options {
		optfn(&mbus)
	}

	return &mbus
}

// WithRabbitmqConn binds an existing RabbitMQ connection
func WithRabbitmqConn(rconn *amqp091.Connection) OptionRabbitmqFunc {
	return func(r *RabbitmqBus) {
		r.conn = rconn
	}
}

// Connect creates a NATS connection
func (r *RabbitmqBus) Connect() error {
	if r.isClosed {
		r.isClosed = false
	}

	if err := r.connect(); err != nil {
		return err
	}

	// connection watcher
	go func() {
		for {
			reason, ok := <-r.conn.NotifyClose(make(chan *amqp091.Error))
			if !ok {
				if r.isClosed {
					return
				}
				r.Log.Error("connection unexpected closed", "reason", reason)

				r.mu.Lock()
				for {
					if connErr := r.connect(); connErr != nil {
						r.Log.Error("connection failed, trying to reconnect", "err", connErr)
						time.Sleep(DefaultReconnectionWait)
						continue
					}
					break
				}
				r.mu.Unlock()
			}
		}
	}()

	return nil
}

// SubscribePing subscribe ping messages
func (r *RabbitmqBus) SubscribePing(topic string, callback PingHandler) (Subscription, error) {

	sub := RmqSubscription{
		Topics:       []string{topic},
		Queue:        exchangePing + "-" + rid.New(ridQueue),
		Exchange:     exchangePing,
		ExchangeKind: amqp091.ExchangeFanout,
		QueueArgs:    amqp091.Table{"x-expires": DefaultQueueExpire},
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}
				callback()
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// SubscribeRequest subscribe request messages
func (r *RabbitmqBus) SubscribeRequest(topic string, callback RequestHandler) (Subscription, error) {

	sub := RmqSubscription{
		Topics:       []string{topic},
		Queue:        topic + "-" + rid.New(ridQueue),
		Exchange:     exchangeRequest,
		ExchangeKind: amqp091.ExchangeTopic,
		QueueArgs:    amqp091.Table{"x-expires": DefaultQueueExpire},
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}

				//callback
				var data proxy.Request
				err = json.Unmarshal(msg.Body, &data)
				if err != nil {
					r.Log.Error("Error unmarshall data", "topic", topic, "error", err)
					continue
				}
				callback(topic, msg.ReplyTo, &data)
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// SubscribeRequests subscribe request messages using multiple topics
func (r *RabbitmqBus) SubscribeRequests(topics []string, callback RequestHandler) (Subscription, error) {

	sub := RmqSubscription{
		Topics:       topics,
		Queue:        rid.New(ridQueue),
		Exchange:     exchangeRequest,
		ExchangeKind: amqp091.ExchangeTopic,
		QueueArgs:    amqp091.Table{"x-expires": DefaultQueueExpire},
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}

				//callback
				var data proxy.Request
				err = json.Unmarshal(msg.Body, &data)
				if err != nil {
					r.Log.Error("Error unmarshall data", "topics", topics, "error", err)
					continue
				}
				callback(msg.RoutingKey, msg.ReplyTo, &data)
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// SubscribeAnnounce subscribe announce messages
func (r *RabbitmqBus) SubscribeAnnounce(topic string, callback AnnounceHandler) (Subscription, error) {

	sub := RmqSubscription{
		Topics:       []string{topic},
		Queue:        topic + "-" + rid.New(ridQueue),
		Exchange:     exchangeAnnounce,
		ExchangeKind: amqp091.ExchangeFanout,
		QueueArgs:    amqp091.Table{"x-expires": DefaultQueueExpire},
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}

				//callback
				var data proxy.Announcement
				err = json.Unmarshal(msg.Body, &data)
				if err != nil {
					r.Log.Error("Error unmarshall data", "topic", topic, "error", err)
					continue
				}
				callback(&data)
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// SubscribeEvent subscribe event messages
func (r *RabbitmqBus) SubscribeEvent(topic string, queue string, callback EventHandler) (Subscription, error) {

	sub := RmqSubscription{
		Topics:       []string{topic},
		Queue:        queue,
		Exchange:     exchangeEvent,
		ExchangeKind: amqp091.ExchangeTopic,
		QueueArgs:    amqp091.Table{"x-message-ttl": DefaultMessageExpire},
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}

				//callback
				callback(msg.Body)
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// SubscribeCreateRequest subscribe create request messages
func (r *RabbitmqBus) SubscribeCreateRequest(topic string, queue string, callback RequestHandler) (Subscription, error) {
	sub := RmqSubscription{
		Topics:       []string{topic},
		Queue:        queue + "-" + topic,
		Exchange:     exchangeRequest,
		ExchangeKind: amqp091.ExchangeTopic,
	}

	d, err := sub.execute(r)
	if err != nil {
		return nil, err
	}
	go func(msgs <-chan amqp091.Delivery) {
		for {
			for msg := range msgs {
				err := msg.Ack(false)
				if err != nil {
					r.Log.Error("failed to ack message: %w", err)
					continue
				}

				//callback
				var data proxy.Request
				err = json.Unmarshal(msg.Body, &data)
				if err != nil {
					r.Log.Error("Error unmarshall data", "topic", topic, "error", err)
					continue
				}
				callback(topic, msg.ReplyTo, &data)
			}
			if r.isClosed {
				return
			}
			msgs = sub.reconnect(r)
		}
	}(d)

	return &sub, nil
}

// PublishResponse sends response message
func (r *RabbitmqBus) PublishResponse(topic string, msg *proxy.Response) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	//exchange should be empty
	return r.publish(topic, "", data)
}

// PublishPing sends ping message
func (r *RabbitmqBus) PublishPing(topic string) error {
	data, err := json.Marshal(&proxy.Request{})
	if err != nil {
		return err
	}
	return r.publish(topic, topic, data)
}

// PublishAnnounce sends announce message
func (r *RabbitmqBus) PublishAnnounce(topic string, msg *proxy.Announcement) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.publish(topic, topic, data)

}

// PublishEvent sends event message
func (r *RabbitmqBus) PublishEvent(topic string, msg ari.Event) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.publish(topic, exchangeEvent, data)

}

// Close closes the connection
func (r *RabbitmqBus) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.isClosed = true
	r.conn.Close() // nolint: errcheck
}

// GetWildcardString returns wildcard based on type
func (r *RabbitmqBus) GetWildcardString(w WildcardType) string {
	switch w {
	case WildcardOneWord:
		return "*"
	case WildcardZeroOrMoreWords:
		return "#"
	}
	return ""
}

// Request sends a request message
func (r *RabbitmqBus) Request(topic string, req *proxy.Request) (*proxy.Response, error) {
	var resp proxy.Response

	requestData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	channel, err := r.conn.Channel()
	r.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	consumerID := rid.New(ridConsumerReq)
	msgs, err := channel.Consume(
		"amq.rabbitmq.reply-to", // queue
		consumerID,              // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return nil, eris.Wrap(err, "Error consumming channel")
	}
	defer channel.Cancel(consumerID, false) // nolint: errcheck

	ctx, cancel := context.WithTimeout(context.Background(), r.Config.RequestTimeout)
	defer cancel()
	for i := 0; i <= r.Config.TimeoutRetries; i++ {
		err = channel.PublishWithContext(
			ctx,
			exchangeRequest, // exchange
			topic,           // routing key
			false,           // mandatory
			false,           // immediate
			amqp091.Publishing{
				ContentType:   "application/json",
				CorrelationId: rid.New(ridCorrelation),
				Body:          requestData,
				ReplyTo:       "amq.rabbitmq.reply-to",
			})

		if errors.Is(err, context.DeadlineExceeded) {
			r.countTimeouts++
			continue
		}
	}

	if err != nil {
		return nil, eris.Wrap(err, "Failed to publish message")
	}

	msg := <-msgs
	if err := json.Unmarshal(msg.Body, &resp); err != nil {
		r.Log.Error("Error on Unmarshal response", "topic", topic, "error", err)
		return nil, err
	}
	return &resp, nil

}

// MultipleRequest sends a request message to multiple consumers
func (r *RabbitmqBus) MultipleRequest(topic string, req *proxy.Request, expectedResp int) ([]*proxy.Response, error) {

	responses := make([]*proxy.Response, 0, expectedResp)

	requestData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	channel, err := r.conn.Channel()
	r.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	consumerID := rid.New(ridConsumerReq)
	msgs, err := channel.Consume(
		"amq.rabbitmq.reply-to", // queue
		consumerID,              // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return nil, eris.Wrap(err, "Error consumming channel")
	}
	defer channel.Cancel(consumerID, false) // nolint: errcheck

	ctx, cancel := context.WithTimeout(context.Background(), r.Config.RequestTimeout)
	defer cancel()
	for i := 0; i <= r.Config.TimeoutRetries; i++ {
		err = channel.PublishWithContext(
			ctx,
			exchangeRequest, // exchange
			topic,           // routing key
			false,           // mandatory
			false,           // immediate
			amqp091.Publishing{
				ContentType:   "application/json",
				CorrelationId: rid.New(ridCorrelation),
				Body:          requestData,
				ReplyTo:       "amq.rabbitmq.reply-to",
			})

		if errors.Is(err, context.DeadlineExceeded) {
			r.countTimeouts++
			continue
		}
	}
	if err != nil {
		return nil, eris.Wrap(err, "Failed to publish message")
	}

	timer := time.NewTimer(r.Config.RequestTimeout)
	defer timer.Stop()
	responseCount := 0
	for {
		select {
		case <-timer.C:
			return responses, nil
		case msg, more := <-msgs:
			if !more {
				return responses, nil
			}
			var resp proxy.Response
			if err := json.Unmarshal(msg.Body, &resp); err != nil {
				r.Log.Error("Error on Unmarshal response", "topic", topic, "error", err)
				return nil, err
			}
			responses = append(responses, &resp)
			responseCount++
			if responseCount >= expectedResp {
				return responses, nil
			}
		}
	}
}

// MultipleRequestReturnFirstGoodResponse sends a request message to multiple consumers and returns the first good response
func (r *RabbitmqBus) MultipleRequestReturnFirstGoodResponse(topic string, req *proxy.Request, expectedResp int) (*proxy.Response, error) {

	requestData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	channel, err := r.conn.Channel()
	r.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	consumerID := rid.New(ridConsumerReq)
	msgs, err := channel.Consume(
		"amq.rabbitmq.reply-to", // queue
		consumerID,              // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return nil, eris.Wrap(err, "Error consumming channel")
	}
	defer channel.Cancel(consumerID, false) // nolint: errcheck

	ctx, cancel := context.WithTimeout(context.Background(), r.Config.RequestTimeout)
	defer cancel()
	for i := 0; i <= r.Config.TimeoutRetries; i++ {
		err = channel.PublishWithContext(
			ctx,
			exchangeRequest, // exchange
			topic,           // routing key
			false,           // mandatory
			false,           // immediate
			amqp091.Publishing{
				ContentType:   "application/json",
				CorrelationId: rid.New(ridCorrelation),
				Body:          requestData,
				ReplyTo:       "amq.rabbitmq.reply-to",
			})

		if errors.Is(err, context.DeadlineExceeded) {
			r.countTimeouts++
			continue
		}
	}
	if err != nil {
		return nil, eris.Wrap(err, "Failed to publish message")
	}

	timer := time.NewTimer(r.Config.RequestTimeout)
	defer timer.Stop()
	responseCount := 0
	for {
		select {
		case <-timer.C:
			// Return the last error if we got one; otherwise, return a timeout error
			if err == nil {
				err = eris.New("timeout")
			}

			return nil, err
		case msg, more := <-msgs:

			if !more {
				if err == nil {
					err = eris.New("no data")
				}
				return nil, err
			}

			var resp proxy.Response
			if err = json.Unmarshal(msg.Body, &resp); err != nil {
				r.Log.Error("Error on Unmarshal response", "topic", topic, "error", err)
				continue
			}

			if err = resp.Err(); err == nil { // store the error for later return
				return &resp, nil // No error means to return the current value
			}

			responseCount++
			if responseCount > expectedResp {
				return nil, eris.New("no data")
			}
		}
	}
}

// TimeoutCount is the amount of times the communication times out
func (r *RabbitmqBus) TimeoutCount() int64 {
	return r.countTimeouts
}

func (r *RabbitmqBus) connect() error {
	var err error
	if r.conn, err = amqp091.Dial(r.Config.URL); err != nil {
		return eris.Wrap(err, "connect")
	}
	if r.channel, err = r.conn.Channel(); err != nil {
		return eris.Wrap(err, "channel create")
	}
	return nil
}

func (r *RabbitmqBus) publish(topic string, exchange string, data []byte) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.channel.PublishWithContext(
		context.Background(),
		exchange, // exchange, for now always using the default exchange
		topic,
		false,
		false,
		amqp091.Publishing{
			Headers:         amqp091.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            data,
			DeliveryMode:    amqp091.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,                 // 0-9
		})
}

func (r *RabbitmqBus) newChannel() (*amqp091.Channel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.conn.Channel()
}
