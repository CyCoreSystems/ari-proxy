package client

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/CyCoreSystems/ari-proxy/v5/client/bus"
	"github.com/CyCoreSystems/ari-proxy/v5/client/cluster"
	"github.com/CyCoreSystems/ari-proxy/v5/messagebus"
	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rotisserie/eris"

	"github.com/inconshreveable/log15"
	"github.com/nats-io/nats.go"
)

// ClosureGracePeriod is the amount of time to wait after the closure of the
// context to close the client.  The delay here is important in order to allow
// wrap-up code to finish processing before losing connection to ARI.
// Depending on the characteristics of traffic and code, this value may need to
// be tweaked.
//
// NOTE: It is always possible to call `Close()` on the client to close it
// immediately.
var ClosureGracePeriod = 10 * time.Second

// DefaultRequestTimeout is the default timeout for a MessageBus request.  (Note: Answer() takes longer than 250ms on average)
var DefaultRequestTimeout = 500 * time.Millisecond

// DefaultInputBufferLength is the default size of the event buffer for events
// coming in from MessageBus
const DefaultInputBufferLength = 100

// DefaultClusterMaxAge is the default maximum age for cluster members to be
// considered by this client
var DefaultClusterMaxAge = 5 * time.Minute

// ErrNil indicates that the request returned an empty response
var ErrNil = eris.New("Nil")

// core is the core, functional piece of a Client which is the same across the
// family of derived clients.  It manages stateful elements such as the bus,
// the MessageBus connection, and the cluster membership
type core struct {
	// cluster describes the cluster of ARI proxies
	cluster *cluster.Cluster

	// clusterMaxAge is the maximum age of cluster members to include in queries
	clusterMaxAge time.Duration

	// inputBufferLength is the size of the buffer for events coming in from MessageBus
	inputBufferLength int

	log log15.Logger

	// Message Bus over which messages will be transceived
	// One of mbus or uri must be specified.
	mbus messagebus.Client

	// prefix is the prefix to use on all MessageBus subjects.  It defaults to "ari.".
	prefix string

	// refCounter is the reference counter for derived clients.  When there are
	// no more referenced clients, the core is shut down.
	refCounter int

	// requestTimeout is the timeout duration of a request
	requestTimeout time.Duration

	// timeoutRetries is the amount of times to retry on message bus timeout
	timeoutRetries int

	// uri provies the URI to which a Message Bus connection should be established. One
	// of mbus or uri must be specified. This option may also be supplied by
	// the `MESSAGEBUS_URL` environment variable.
	uri string

	// annSub is the subscription to proxy announcements
	annSub messagebus.Subscription

	// closeChan is the signal channel responsible for shutting down core
	// services.  When it is closed, all core services should exit.
	closeChan chan struct{}

	// closed indicates the core has been closed
	closed bool

	// closeMBusOnClose indicates that the Message Bus connection should be closed when
	// the ari.Client is closed
	closeMBusOnClose bool

	// started indicates whether this core has been started; a started core will
	// no-op core.start()
	started bool
}

// clientClosed is called any time a derived ARI client is closed; if the
// reference counter is ever dropped to zero, the core is also shut down
func (c *core) ClientClosed() {
	c.refCounter--

	if c.refCounter < 1 {
		c.close()
	}
}

// close shuts down the core
func (c *core) close() {
	if !c.closed {
		c.closed = true
		close(c.closeChan)
	}

	if c.annSub != nil {
		err := c.annSub.Unsubscribe()
		if err != nil {
			c.log.Debug("failed to unsubscribe from proxy announcements", "error", err)
		}
	}

	c.mbus.Close()
}

func (c *core) Start() error {
	// increment the client reference counter
	c.refCounter++

	// Only start the core once
	if c.started {
		return nil
	}
	c.started = true

	c.closeChan = make(chan struct{})

	// Connect to MessageBus, if we do not already have a connection
	if c.mbus == nil {
		switch messagebus.GetType(c.uri) {
		case messagebus.TypeNats:
			c.mbus = &messagebus.NatsBus{
				Config: messagebus.Config{
					URL:            c.uri,
					TimeoutRetries: c.timeoutRetries,
					RequestTimeout: c.requestTimeout,
				},
				Log: c.log,
			}
		case messagebus.TypeRabbitmq:
			c.mbus = &messagebus.RabbitmqBus{
				Config: messagebus.Config{
					URL:            "amqp://guest:guest@rabbitmq:5672/",
					TimeoutRetries: c.timeoutRetries,
					RequestTimeout: c.requestTimeout,
				},
				Log: c.log,
			}
		default:
			return errors.New("Unknown url for MessageBus: " + c.uri)
		}

		err := c.mbus.Connect()
		if err != nil {
			c.close()
			return eris.Wrap(err, "failed to connect to MessageBus")
		}
		c.closeMBusOnClose = true
	}

	// Create and start the cluster
	c.cluster = cluster.New()

	// Maintain the cluster
	err := c.maintainCluster()
	if err != nil {
		c.close()
		return eris.Wrap(err, "failed to start cluster maintenance")
	}

	return nil
}

func (c *core) maintainCluster() (err error) {

	c.annSub, err = c.mbus.SubscribeAnnounce(proxy.AnnouncementSubject(c.prefix), func(o *proxy.Announcement) {
		c.cluster.Update(o.Node, o.Application)
	})
	if err != nil {
		return eris.Wrap(err, "failed to listen to proxy announcements")
	}

	// Send an initial ping for proxy announcements
	err = c.mbus.PublishPing(proxy.PingSubject(c.prefix))
	if err != nil {
		return eris.Wrap(err, "failed to publish ping")
	}
	return err
}

// Client provides an ari.Client for an ari-proxy server
type Client struct {
	*core

	bus ari.Bus

	appName string

	cancel context.CancelFunc

	// closed indicates that this client has been closed and is no longer attached to a core
	closed bool
}

// New creates a new Client to the Asterisk ARI NATS/RabbitMQ proxy.
func New(ctx context.Context, opts ...OptionFunc) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &Client{
		appName: os.Getenv("ARI_APPLICATION"),
		core: &core{
			cluster:           cluster.New(),
			clusterMaxAge:     DefaultClusterMaxAge,
			inputBufferLength: DefaultInputBufferLength,
			log:               log15.New(),
			prefix:            "ari.",
			requestTimeout:    DefaultRequestTimeout,
			uri:               "nats://localhost:4222",
		},
		cancel: cancel,
	}
	c.log.SetHandler(log15.DiscardHandler())

	// Load environment-based configurations
	if os.Getenv("MESSAGEBUS_URL") != "" {
		c.core.uri = os.Getenv("MESSAGEBUS_URL")
	} else if os.Getenv("NATS_URI") != "" { //backward compatibility
		c.core.uri = os.Getenv("NATS_URI")
	}

	// Load explicit configurations
	for _, opt := range opts {
		opt(c)
	}

	// Start the core, if it is not already started
	err := c.core.Start()
	if err != nil {
		return nil, eris.Wrap(err, "failed to start core")
	}

	// Create the bus
	c.bus = bus.New(c.core.prefix, c.core.mbus, c.core.log)

	// Call Close whenever the context is closed
	go func() {
		<-ctx.Done()
		if !c.closed {
			// Only wait the grace period if we have not
			// already been closed.
			<-time.After(ClosureGracePeriod)
		}

		c.Close()
	}()

	return c, nil
}

// New returns a new client from the existing one.  The new client will have a
// separate event bus and lifecycle, allowing the closure of all subscriptions
// and handles derived from the client by simply closing the client.  The
// underlying MessageBus connection and cluster awareness (the common Core) will be
// preserved across derived Client lifecycles.
func (c *Client) New(ctx context.Context) *Client {
	_, cancel := context.WithCancel(ctx)

	return &Client{
		appName: c.appName,
		cancel:  cancel,
		core:    c.core,
		bus:     bus.New(c.core.prefix, c.core.mbus, c.core.log),
	}
}

// OptionFunc is a function which configures options on a Client
type OptionFunc func(*Client)

// FromClient configures the ARI Application to use the transport details from
// another ARI Client.  Transport-related details are copied, such as the MessageBus
// Client, the MessageBus prefix, the timeout values.
//
// Specifically NOT copied are dialog, application, and asterisk details.
//
// NOTE: use of this function will cause MessageBus connection leakage if there is a
// mix of uses of FromClient and not over a period of time.  If you intend to
// use FromClient, it is recommended that you always pass a MessageBus client in to
// the first ari.Client and maintain lifecycle control of it manually.
func FromClient(cl ari.Client) OptionFunc {
	return func(c *Client) {
		old, ok := cl.(*Client)
		if ok {
			c.core = old.core
		}
	}
}

// WithApplication configures the ARI Client to use the provided ARI Application
func WithApplication(name string) OptionFunc {
	return func(c *Client) {
		c.appName = name
	}
}

// WithLogger sets the logger on a Client.
func WithLogger(l log15.Logger) OptionFunc {
	return func(c *Client) {
		c.log = l
	}
}

// WithLogHandler sets the logging handler on a Client logger
func WithLogHandler(h log15.Handler) OptionFunc {
	return func(c *Client) {
		c.log.SetHandler(h)
	}
}

// WithURI sets the MessageBus URI to which the client will attempt to connect.
// The MessageBus URI may also be configured by the environment variable `MESSAGEBUS_URL`.
func WithURI(uri string) OptionFunc {
	return func(c *Client) {
		c.core.uri = uri
	}
}

// WithNATS binds an existing NATS connection
func WithNATS(nc *nats.EncodedConn) OptionFunc {
	return func(c *Client) {
		c.mbus = messagebus.NewNatsBus(
			messagebus.Config{},
			messagebus.WithNatsConn(nc),
		)
	}
}

// WithRabbitmq binds an existing RabbitMQ connection
func WithRabbitmq(conn *amqp091.Connection) OptionFunc {
	return func(c *Client) {
		c.mbus = messagebus.NewRabbitmqBus(
			messagebus.Config{},
			messagebus.WithRabbitmqConn(conn),
		)
	}
}

// WithPrefix configures the MessageBus Prefix to use on a Client
func WithPrefix(prefix string) OptionFunc {
	return func(c *Client) {
		c.core.prefix = prefix
	}
}

// WithTimeoutRetries configures the amount of times to retry on request timeout for a Client
func WithTimeoutRetries(count int) OptionFunc {
	return func(c *Client) {
		c.core.timeoutRetries = count
	}
}

// ApplicationName returns the ARI application's name
func (c *Client) ApplicationName() string {
	return c.appName
}

// Connected indicates whether the client is connected through to at least one ARI websocket
func (c *Client) Connected() bool {
	if c.closed {
		return false
	}

	// TODO: this is a surrogate indicator with low resolution... we should have
	// something more proactive and concrete
	if len(c.cluster.App(c.appName, proxy.AnnouncementInterval*2)) < 1 {
		return false
	}
	return true
}

// Close shuts down the client
func (c *Client) Close() {
	if c.cancel != nil {
		c.cancel()
	}

	if c.bus != nil {
		c.bus.Close()
	}

	if !c.closed && c.core != nil {
		c.closed = true
		c.core.ClientClosed()
	}
}

// Application is the application operation accessor
func (c *Client) Application() ari.Application {
	return &application{c}
}

// Asterisk is the asterisk operation accessor
func (c *Client) Asterisk() ari.Asterisk {
	return &asterisk{c}
}

// Bridge is the bridge operation accessor
func (c *Client) Bridge() ari.Bridge {
	return &bridge{c}
}

// Bus is the bus operation accessor
func (c *Client) Bus() ari.Bus {
	return c.bus
}

// Channel is the channel operation accessor
func (c *Client) Channel() ari.Channel {
	return &channel{c}
}

// DeviceState is the device state operation accessor
func (c *Client) DeviceState() ari.DeviceState {
	return &deviceState{c}
}

// Endpoint is the endpoint accessor
func (c *Client) Endpoint() ari.Endpoint {
	return &endpoint{c}
}

// LiveRecording is the live recording accessor
func (c *Client) LiveRecording() ari.LiveRecording {
	return &liveRecording{c}
}

// Mailbox is the mailbox accessor
func (c *Client) Mailbox() ari.Mailbox {
	return &mailbox{c}
}

// Playback is the media playback accessor
func (c *Client) Playback() ari.Playback {
	return &playback{c}
}

// Sound is the sound accessor
func (c *Client) Sound() ari.Sound {
	return &sound{c}
}

// StoredRecording is the stored recording accessor
func (c *Client) StoredRecording() ari.StoredRecording {
	return &storedRecording{c}
}

// TextMessage is the text message accessor
func (c *Client) TextMessage() ari.TextMessage {
	return nil
}

func (c *Client) commandRequest(req *proxy.Request) error {
	resp, err := c.makeRequest("command", req)
	if err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) createRequest(req *proxy.Request) (*ari.Key, error) {
	resp, err := c.makeRequest("create", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Key == nil {
		return nil, ErrNil
	}
	return resp.Key, nil
}

func (c *Client) getRequest(req *proxy.Request) (*ari.Key, error) {
	resp, err := c.makeRequest("get", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Key == nil {
		return nil, ErrNil
	}
	return resp.Key, nil
}

func (c *Client) dataRequest(req *proxy.Request) (*proxy.EntityData, error) {
	resp, err := c.makeRequest("data", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Data == nil {
		return nil, ErrNil
	}
	return resp.Data, nil
}

func (c *Client) listRequest(req *proxy.Request) ([]*ari.Key, error) {
	var list []*ari.Key

	responses, err := c.makeRequests("get", req)
	if err != nil {
		return nil, err
	}

	for _, r := range responses {
		err = r.Err()
		if r.Err() != nil || r.Keys == nil {
			continue
		}
		list = append(list, r.Keys...)
	}
	return list, err
}

func (c *Client) makeRequest(class string, req *proxy.Request) (*proxy.Response, error) {
	if !c.completeCoordinates(req) {
		return c.makeBroadcastRequestReturnFirstGoodResponse(class, req)
	}

	c.log.Error("request", "class", class, "req", req, "subject", c.subject(class, req))
	return c.mbus.Request(c.subject(class, req), req)
}

func (c *Client) makeRequests(class string, req *proxy.Request) (responses []*proxy.Response, err error) {
	if req == nil {
		return nil, eris.New("empty request")
	}
	if req.Key == nil {
		req.Key = ari.NewKey("", "")
	}

	expected := len(c.core.cluster.Matching(req.Key.Node, req.Key.App, c.core.clusterMaxAge))

	return c.mbus.MultipleRequest(c.subject(class, req), req, expected)
}

func (c *Client) makeBroadcastRequestReturnFirstGoodResponse(class string, req *proxy.Request) (*proxy.Response, error) {
	if req == nil {
		return nil, eris.New("empty request")
	}

	if req.Key == nil {
		req.Key = ari.NewKey("", "")
	}

	return c.mbus.MultipleRequestReturnFirstGoodResponse(
		c.subject(class, req),
		req,
		len(c.core.cluster.Matching(req.Key.Node, req.Key.App, c.core.clusterMaxAge)),
	)
}

func (c *Client) completeCoordinates(req *proxy.Request) bool {
	if req == nil || req.Key == nil {
		return false
	}

	// coordinates are complete if we have both app and node
	return req.Key.App != "" && req.Key.Node != ""
}

func (c *Client) subject(class string, req *proxy.Request) string {
	if req == nil || req.Key == nil {
		return proxy.Subject(c.core.prefix, class, c.appName, "")
	}
	return proxy.Subject(c.core.prefix, class, req.Key.App, req.Key.Node)
}

// TimeoutCount is the amount of times the communication times out
func (c *Client) TimeoutCount() int64 {
	// countTimeouts tracks how many timeouts the client has received, for metrics.
	return c.mbus.TimeoutCount()
}
