package client

import (
	"context"
	"fmt"
	"os"
	"time"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client/bus"
	"github.com/CyCoreSystems/ari-proxy/client/cluster"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// DefaultInputBufferLength is the default size of the event buffer for events coming in from NATS
const DefaultInputBufferLength = 100

// DefaultClusterMaxAge is the default maximum age for cluster members to be considered by this client
var DefaultClusterMaxAge = 5 * time.Minute

// ErrNil indicates that the request returned an empty response
var ErrNil = errors.New("Nil")

// core is the core, functional piece of a Client which is the same across the family of derived clients.  It manages stateful elements such as the bus, the NATS connection, and the cluster membership
type core struct {
	// cluster describes the cluster of ARI proxies
	cluster *cluster.Cluster

	// clusterMaxAge is the maximum age of cluster members to include in queries
	clusterMaxAge time.Duration

	// inputBufferLength is the size of the buffer for events coming in from NATS
	inputBufferLength int

	log log15.Logger

	// nc provides the nats.EncodedConn over which messages will be transceived.  One of NATS or NATSURI must be specified.
	nc *nats.EncodedConn

	// prefix is the prefix to use on all NATS subjects.  It defaults to "ari.".
	prefix string

	// readOperationRetryCount is the amount of times to retry a read operation
	readOperationRetryCount int

	// refCounter is the reference counter for derived clients.  When there are no more referenced clients, the core is shut down.
	refCounter int

	// requestTimeout is the timeout duration of a request
	requestTimeout time.Duration

	// timeoutRetries is the amount of times to retry on nats timeout
	timeoutRetries int

	// countTimeouts tracks how many timeouts the client has received, for metrics.
	countTimeouts int64

	// uri provies the URI to which a NATS connection should be established. One of NATS or NATSURI must be specified. This option may also be supplied by the `NATS_URI` environment variable.
	uri string

	// annSub is the NATS subscription to proxy announcements
	annSub *nats.Subscription

	// closeChan is the signal channel responsible for shutting down core services.  When it is closed, all core services should exit.
	closeChan chan struct{}

	// closed indicates the core has been closed
	closed bool

	// closeNATSOnClose indicates that the NATS connection should be closed when the ari.Client is closed
	closeNATSOnClose bool

	// started indicates whether this core has been started; a started core will no-op core.start()
	started bool
}

// clientClosed is called any time a derived ARI client is closed; if the reference counter is ever dropped to zero, the core is also shut down
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
			c.log.Debug("failed to unsubscribe from NATS proxy announcements", "error", err)
		}
	}

	if c.closeNATSOnClose && c.nc != nil {
		c.nc.Close()
	}
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

	// Connect to NATS, if we do not already have a connection
	if c.nc == nil {
		n, err := nats.Connect(c.uri)
		if err != nil {
			c.close()
			return errors.Wrap(err, "failed to connect to NATS")
		}

		c.nc, err = nats.NewEncodedConn(n, nats.JSON_ENCODER)
		if err != nil {
			n.Close() // need this here because nc is not yet bound to the core
			c.close()
			return errors.Wrap(err, "failed to encode NATS connection")
		}

		c.closeNATSOnClose = true
	}

	// Create and start the cluster
	c.cluster = cluster.New()

	// Maintain the cluster
	err := c.maintainCluster()
	if err != nil {
		c.close()
		return errors.Wrap(err, "failed to start cluster maintenance")
	}

	return nil
}

func (c *core) maintainCluster() (err error) {
	c.annSub, err = c.nc.Subscribe(proxy.AnnouncementSubject(c.prefix), func(o *proxy.Announcement) {
		c.cluster.Update(o.Node, o.Application)
	})
	if err != nil {
		return errors.Wrap(err, "failed to listen to proxy announcements")
	}

	// Send an initial ping for proxy announcements
	return c.nc.Publish(proxy.PingSubject(c.prefix), &proxy.Request{})
}

// Client provides an ari.Client for an ari-proxy server
type Client struct {
	*core

	bus ari.Bus

	metadata *proxy.Metadata

	cancel context.CancelFunc

	// closed indicates that this client has been closed and is no longer attached to a core
	closed bool
}

// New creates a new Client to the Asterisk ARI NATS proxy.
func New(ctx context.Context, opts ...OptionFunc) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &Client{
		core: &core{
			cluster:           cluster.New(),
			clusterMaxAge:     DefaultClusterMaxAge,
			inputBufferLength: DefaultInputBufferLength,
			log:               log15.New(),
			prefix:            "ari.",
			requestTimeout:    DefaultRequestTimeout,
			uri:               "nats://localhost:4222",
		},
		metadata: &proxy.Metadata{
			Application: os.Getenv("ARI_APPLICATION"),
		},
		cancel: cancel,
	}
	c.log.SetHandler(log15.DiscardHandler())

	// Load environment-based configurations
	if os.Getenv("NATS_URI") != "" {
		c.core.uri = os.Getenv("NATS_URI")
	}

	// Load explicit configurations
	for _, opt := range opts {
		opt(c)
	}

	// Create the bus
	c.bus = &bus.Bus{
		nc:     c.core.nc,
		log:    c.core.log,
		prefix: c.core.prefix,
	}

	// Start the core, if it is not already started
	c.core.Start()

	// Bind events from NATS to our bus
	go c.bindEvents(ctx)

	// Call Close whenever the context is closed
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	return c, nil
}

// New returns a new client from the existing one.  The new client will have a
// separate event bus and lifecycle, allowing the closure of all subscriptions
// and handles derived from the client by simply closing the client.  The
// underlying NATS connection and cluster awareness (the common Core) will be
// preserved across derived Client lifecycles.
func (c *Client) New() *Client {
	return &Client{
		core: c.core,
		bus: &bus.Bus{
			nc:     c.core.nc,
			log:    c.core.log,
			prefix: c.core.prefix,
		},
	}
}

// OptionFunc is a function which configures options on a Client
type OptionFunc func(*Client)

// FromClient configures the ARI Application to use the transport details from
// another ARI Client.  Transport-related details are copied, such as the NATS
// Client, the NATS prefix, the timeout values.
//
// Specifically NOT copied are dialog, application, and asterisk details.
//
// NOTE: use of this function will cause NATS connection leakage if there is a
// mix of uses of FromClient and not over a period of time.  If you intend to
// use FromClient, it is recommended that you always pass a NATS client in to
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
		c.metadata.Application = name
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

// WithNATS binds an existing NATS connection
func WithNATS(nc *nats.EncodedConn) OptionFunc {
	return func(c *Client) {
		c.nc = nc
	}
}

// WithPrefix configures the NATS Prefix to use on a Client
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
	if c.metadata == nil {
		return ""
	}
	return c.metadata.Application
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
	var resp *proxy.Response
	var err error

	// if we have complete coordinates, we can make a direct request
	if c.completeCoordinates(req) {
		resp, err = c.makeRequest("command", req)
	} else {
		resp, err = c.makeBroadcastRequestReturnFirstGoodResponse("command", req)
	}

	if err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) createRequest(req *proxy.Request) (*proxy.Entity, error) {
	resp, err := c.makeRequest("create", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Entity == nil {
		return nil, ErrNil
	}
	return resp.Entity, nil
}

func (c *Client) getRequest(req *proxy.Request) (*proxy.Entity, error) {
	// If we do not have complete coordinates, we cannot make a reasonable get request
	if !c.completeCoordinates(req) {
		return nil, errors.New("Incomplete coordinates")
	}

	resp, err := c.makeRequest("get", req)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.Entity == nil {
		return nil, ErrNil
	}
	return resp.Entity, nil
}

func (c *Client) dataRequest(req *proxy.Request) (*proxy.EntityData, error) {
	var resp *proxy.Response
	var err error

	// if we have complete coordinates, we can make a direct request
	if c.completeCoordinates(req) {
		resp, err = c.makeRequest("data", req)
	} else {
		resp, err = c.makeBroadcastRequestReturnFirstGoodResponse("data", req)
	}

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

func (c *Client) listRequest(req *proxy.Request) (*proxy.EntityList, error) {
	var list proxy.EntityList

	responses, err := c.makeBroadcastRequest("get", req)
	if err != nil {
		return nil, err
	}

	for _, r := range responses {
		if r.Err() != nil || r.EntityList == nil {
			continue
		}
		list.List = append(list.List, r.EntityList.List...)
	}
	return &list, nil
}

func (c *Client) makeRequest(class string, req *proxy.Request) (*proxy.Response, error) {
	var resp proxy.Response
	var err error
	var retries int

	for retries <= c.core.timeoutRetries {
		retries++
		err = c.nc.Request(c.subject(class, req), req, &resp, c.requestTimeout)
		if err == nats.ErrTimeout {
			continue
		}
		return &resp, err
	}

	return nil, err
}

func (c *Client) makeBroadcastRequest(class string, req *proxy.Request) ([]*proxy.Response, error) {
	var responses []*proxy.Response
	var err error

	var responseCount int
	expected := len(c.core.cluster.Matching(c.nodeForRequest(req), c.appForRequest(req), c.core.clusterMaxAge))
	reply := uuid.NewV1().String()
	replyChan := make(chan *proxy.Response)
	replySub, err := c.core.nc.Subscribe(reply, func(o *proxy.Response) {
		responseCount++

		replyChan <- o

		if responseCount >= expected {
			close(replyChan)
		}
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to data responses")
	}
	defer replySub.Unsubscribe()

	// Make an all-call for the entity data
	err = c.core.nc.PublishRequest(c.subject(class, req), reply, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request for data")
	}

	// Wait for replies
	for {
		select {
		case <-time.After(c.requestTimeout):
			return responses, nil
		case resp, ok := <-replyChan:
			if !ok {
				return responses, nil
			}
			responses = append(responses, resp)
		}
	}
}

func (c *Client) makeBroadcastRequestReturnFirstGoodResponse(class string, req *proxy.Request) (*proxy.Response, error) {
	var responseCount int
	expected := len(c.core.cluster.Matching(c.nodeForRequest(req), c.appForRequest(req), c.core.clusterMaxAge))
	reply := uuid.NewV1().String()
	replyChan := make(chan *proxy.Response)
	replySub, err := c.core.nc.Subscribe(reply, func(o *proxy.Response) {
		responseCount++

		if o.Err() == nil {
			replyChan <- o
		}

		if responseCount >= expected {
			close(replyChan)
		}
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to data responses")
	}
	defer replySub.Unsubscribe()

	// Make an all-call for the entity data
	err = c.core.nc.PublishRequest(c.subject(class, req), reply, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request for data")
	}

	// Wait for replies
	select {
	case <-time.After(c.requestTimeout):
		return nil, errors.New("timeout")
	case resp := <-replyChan:
		if resp == nil {
			return nil, proxy.ErrNotFound
		}
		return resp, nil
	}
}

func (c *Client) appForRequest(req *proxy.Request) string {
	var app string
	if c.metadata != nil {
		app = c.metadata.Application
	}
	if req.Metadata != nil && req.Metadata.Application != "" {
		app = req.Metadata.Application
	}
	return app
}

func (c *Client) nodeForRequest(req *proxy.Request) string {
	var node string
	if c.metadata != nil {
		node = c.metadata.Node
	}
	if req.Metadata != nil && req.Metadata.Node != "" {
		node = req.Metadata.Node
	}
	return node
}

func (c *Client) completeCoordinates(req *proxy.Request) bool {
	// coordinates are complete if we have both app and node
	return c.appForRequest(req) != "" &&
		c.nodeForRequest(req) != ""
}

func (c *Client) subject(class string, req *proxy.Request) string {
	return proxy.Subject(c.core.prefix, class, c.appForRequest(req), c.nodeForRequest(req))
}

func (c *Client) eventsSubject() (subj string) {
	// Attempt to build events subject from metadata
	if c.metadata != nil {
		if c.metadata.Application != "" {
			subj = fmt.Sprintf("%sevent.%s", c.core.prefix, c.metadata.Application)

			if c.metadata.Node != "" {
				subj += "." + c.metadata.Node
			} else {
				subj += ".>"
			}
		}

		// a dialog always overrides
		if c.metadata.Dialog != "" {
			subj = fmt.Sprintf("%sdialogevent.%s", c.core.prefix, c.metadata.Dialog)
		}
	}

	// If we still have no subject, listen to everything
	if subj == "" {
		subj = fmt.Sprintf("%sevent.>", c.core.prefix)
	}
	return
}
