package messagebus

import (
	"sync"
	"time"

	"github.com/CyCoreSystems/ari/v5/rid"
	"github.com/rabbitmq/amqp091-go"
)

// RmqSubscription handle RabbitMQ subscription
type RmqSubscription struct {
	consumerID string
	channel    *amqp091.Channel
	mu         sync.RWMutex

	Topics       []string
	Queue        string
	Exchange     string
	ExchangeKind string
	QueueArgs    amqp091.Table
}

// Unsubscribe remove the subscription
func (rs *RmqSubscription) Unsubscribe() error {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	if rs.consumerID == "" {
		return nil
	}
	err := rs.channel.Cancel(rs.consumerID, false)
	if err != nil {
		return err
	}
	return rs.channel.Close()
}

// reconnect reconnects the subscription
func (rs *RmqSubscription) reconnect(r *RabbitmqBus) <-chan amqp091.Delivery {
	rs.mu.Lock()

	if rs.channel != nil {
		rs.channel.Close() // nolint: errcheck
	}
	rs.mu.Unlock()

	for {
		msgs, err := rs.execute(r)
		if err != nil {
			r.Log.Error("failed to execute subscription", "error", err)
			time.Sleep(DefaultReconnectionWait)
			continue
		}
		return msgs
	}
}

// execute declares the subscription on RabbitMQ
func (rs *RmqSubscription) execute(r *RabbitmqBus) (<-chan amqp091.Delivery, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	ch, err := r.newChannel()
	if err != nil {
		return nil, err
	}

	rs.channel = ch
	queue, err := ch.QueueDeclare(
		rs.Queue,     // name of queue
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // nowait
		rs.QueueArgs, // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		rs.Exchange,     // name of exchange
		rs.ExchangeKind, // kind
		true,            // durable
		false,           // delete when unused
		false,           // internal
		false,           // nowait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}
	if queue.Name != rs.Queue {
		rs.Queue = queue.Name
	}

	for _, topic := range rs.Topics {
		err = ch.QueueBind(queue.Name, topic, rs.Exchange, false, nil)
		if err != nil {
			return nil, err
		}
	}

	rs.consumerID = rid.New(ridConsumer)
	return ch.Consume(rs.Queue, rs.consumerID, false, false, true, true, nil)
}
