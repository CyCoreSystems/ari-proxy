package messagebus

import (
	"time"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/rid"
	"github.com/inconshreveable/log15"
	"github.com/nats-io/nats.go"
	"github.com/rotisserie/eris"
)

// NatsBus is MessageBus implementation for RabbitMQ
type NatsBus struct {
	Config Config
	Log    log15.Logger

	conn          *nats.EncodedConn
	countTimeouts int64
}

// OptionNatsFunc options for RabbitMQ
type OptionNatsFunc func(n *NatsBus)

// NatsMSubscription handle multiple subscriptions with same handler
type NatsMSubscription struct {
	Subscriptions []*nats.Subscription
}

// Unsubscribe removes the multiple subscriptions
func (n *NatsMSubscription) Unsubscribe() error {
	for _, sub := range n.Subscriptions {
		err := sub.Unsubscribe()
		if err != nil {
			return eris.Wrap(err, "failed to unsubscribe "+sub.Subject)
		}
	}
	return nil
}

// NewNatsBus creates a NatsBus
func NewNatsBus(config Config, options ...OptionNatsFunc) *NatsBus {

	mbus := NatsBus{
		Config: config,
	}

	for _, optfn := range options {
		optfn(&mbus)
	}

	return &mbus
}

// WithNatsConn binds an existing NATS connection
func WithNatsConn(nconn *nats.EncodedConn) OptionNatsFunc {
	return func(n *NatsBus) {
		n.conn = nconn
	}
}

// Connect creates a NATS connection
func (n *NatsBus) Connect() error {
	reconnectionAttempts := DefaultReconnectionAttemts
	nc, err := nats.Connect(n.Config.URL,
		nats.Name(n.Config.ID),
		nats.MaxReconnects(DefaultReconnectionAttemts),
		nats.ReconnectWait(DefaultReconnectionWait),
		nats.ReconnectHandler(func(c *nats.Conn) {
			reconnectionAttempts--

			n.Log.Info("retrying to connect to NATS server", "attempts", reconnectionAttempts)
		}),
		nats.MaxPingsOutstanding(3),
	)
	if err != nil {
		return eris.Wrap(err, "failed to connect to NATS")
	}
	/*

		for err == nats.ErrNoServers && reconnectionAttempts > 0 {
			n.Log.Info("retrying to connect to NATS server", "attempts", reconnectionAttempts)
			time.Sleep(DefaultReconnectionWait)
			nc, err = nats.Connect(n.Config.URL)
			reconnectionAttempts--
		}
	*/
	n.conn, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		nc.Close()
		return eris.Wrap(err, "failed to encode NATS connection")
	}
	return nil
}

// SubscribePing subscribe ping messages
func (n *NatsBus) SubscribePing(topic string, callback PingHandler) (Subscription, error) {
	return n.conn.Subscribe(topic, func(m *nats.Msg) {
		callback()
	})
}

// SubscribeRequest subscribe request messages
func (n *NatsBus) SubscribeRequest(topic string, callback RequestHandler) (Subscription, error) {
	return n.conn.Subscribe(topic, callback)
}

// SubscribeRequests subscribe request messages using multiple topics
func (n *NatsBus) SubscribeRequests(topics []string, callback RequestHandler) (Subscription, error) {

	subs := NatsMSubscription{}
	for _, topic := range topics {
		sub, err := n.conn.Subscribe(topic, callback)
		if err != nil {
			subs.Unsubscribe() // nolint: errcheck
			return nil, eris.Wrapf(err, "failed to create %s subscription", topic)
		}
		subs.Subscriptions = append(subs.Subscriptions, sub)

	}
	return &subs, nil
}

// SubscribeAnnounce subscribe announce messages
func (n *NatsBus) SubscribeAnnounce(topic string, callback AnnounceHandler) (Subscription, error) {
	return n.conn.Subscribe(topic, callback)
}

// SubscribeEvent subscribe event messages
func (n *NatsBus) SubscribeEvent(topic string, queue string, callback EventHandler) (Subscription, error) {
	return n.conn.Subscribe(topic, func(m *nats.Msg) {
		callback(m.Data)
	})
}

// SubscribeCreateRequest subscribe create request messages
func (n *NatsBus) SubscribeCreateRequest(topic string, queue string, callback RequestHandler) (Subscription, error) {
	return n.conn.QueueSubscribe(topic, queue, callback)
}

// PublishResponse sends response message
func (n *NatsBus) PublishResponse(topic string, msg *proxy.Response) error {
	return n.conn.Publish(topic, msg)
}

// PublishPing sends ping message
func (n *NatsBus) PublishPing(topic string) error {
	return n.conn.Publish(topic, &proxy.Request{})
}

// PublishAnnounce sends announce message
func (n *NatsBus) PublishAnnounce(topic string, msg *proxy.Announcement) error {
	return n.conn.Publish(topic, msg)
}

// PublishEvent sends event message
func (n *NatsBus) PublishEvent(topic string, msg ari.Event) error {
	return n.conn.Publish(topic, msg)
}

// Close closes the connection
func (n *NatsBus) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}

// GetWildcardString returns wildcard based on type
func (n *NatsBus) GetWildcardString(w WildcardType) string {
	switch w {
	case WildcardOneWord:
		return "*"
	case WildcardZeroOrMoreWords:
		return ">"
	}
	return ""
}

// Request sends a request message
func (n *NatsBus) Request(topic string, req *proxy.Request) (*proxy.Response, error) {
	var err error
	var resp proxy.Response
	for i := 0; i <= n.Config.TimeoutRetries; i++ {
		err = n.conn.Request(topic, req, &resp, n.Config.RequestTimeout)
		if err == nats.ErrTimeout {
			n.countTimeouts++
			continue
		}
		if err != nil {
			return nil, err
		}
		return &resp, nil
	}
	return nil, err
}

// MultipleRequest sends a request message to multiple consumers
func (n *NatsBus) MultipleRequest(topic string, req *proxy.Request, expectedResp int) ([]*proxy.Response, error) {
	var responses []*proxy.Response

	reply := rid.New("rp")

	rf := &responseForwarder{
		expected: expectedResp,
		fwdChan:  make(chan *proxy.Response),
	}

	replySub, err := n.conn.Subscribe(reply, rf.Forward)
	if err != nil {
		return nil, eris.Wrap(err, "failed to subscribe to data responses")
	}
	defer replySub.Unsubscribe() // nolint: errcheck

	// Make an all-call for the entity data
	err = n.conn.PublishRequest(topic, reply, req)
	if err != nil {
		return nil, eris.Wrap(err, "failed to make request for data")
	}

	// Wait for replies
	timer := time.NewTimer(n.Config.RequestTimeout)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return responses, nil
		case resp, more := <-rf.fwdChan:
			if !more {
				return responses, nil
			}
			responses = append(responses, resp)
		}
	}
}

// MultipleRequestReturnFirstGoodResponse sends a request message to multiple consumers and returns the first good response
func (n *NatsBus) MultipleRequestReturnFirstGoodResponse(topic string, req *proxy.Request, expectedResp int) (*proxy.Response, error) {

	reply := rid.New("rp")

	rf := &responseForwarder{
		expected: expectedResp,
		fwdChan:  make(chan *proxy.Response),
	}

	replySub, err := n.conn.Subscribe(reply, rf.Forward)
	if err != nil {
		return nil, eris.Wrap(err, "failed to subscribe to data responses")
	}
	defer replySub.Unsubscribe() // nolint: errcheck

	// Make an all-call for the entity data
	err = n.conn.PublishRequest(topic, reply, req)
	if err != nil {
		return nil, eris.Wrap(err, "failed to make request for data")
	}

	// Wait for replies
	timer := time.NewTimer(n.Config.RequestTimeout)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// Return the last error if we got one; otherwise, return a timeout error
			if err == nil {
				err = eris.New("timeout")
			}

			return nil, err
		case resp, more := <-rf.fwdChan:
			if !more {
				if err == nil {
					err = eris.New("no data")
				}

				return nil, err
			}
			if resp != nil {
				if err = resp.Err(); err == nil { // store the error for later return
					return resp, nil // No error means to return the current value
				}
			}
		}
	}
}

// TimeoutCount is the amount of times the communication times out
func (n *NatsBus) TimeoutCount() int64 {
	return n.countTimeouts
}
