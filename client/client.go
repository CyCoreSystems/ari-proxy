package client

import (
	"context"
	"fmt"
	"os"
	"time"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// DefaultInputBufferLength is the default size of the event buffer for events coming in from NATS
const DefaultInputBufferLength = 100

// ErrNil indicates that the request returned an empty response
var ErrNil = errors.New("Nil")

// Client provides an ari.Client for an ari-proxy server
type Client struct {

	// appName provides an ARI application to which to bind.  At least one of Application or Dialog must be specified.  This option may also be supplied by the `ARI_APPLICATION` environment variable.
	appName string

	// asterisk is a unique identifier for a specific Asterisk node.  If specified, all events and all commands will be filtered to/by this Asterisk node.
	asterisk string

	// dialog provides a dialog ID to which to bind.  At least one of Application or Dialog must be specified.
	dialog string

	// inputBufferLength is the size of the buffer for events coming in from NATS
	inputBufferLength int

	// nc provides the nats.EncodedConn over which messages will be transceived.  One of NATS or NATSURI must be specified.
	nc *nats.EncodedConn

	// uri provies the URI to which a NATS connection should be established. One of NATS or NATSURI must be specified. This option may also be supplied by the `NATS_URI` environment variable.
	uri string

	// prefix is the prefix to use on all NATS subjects.  It defaults to "ari.".
	prefix string

	// closeNATSOnClose indicates that the NATS connection should be closed when the ari.Client is closed
	closeNATSOnClose bool

	// readOperationRetryCount is the amount of times to retry a read operation
	readOperationRetryCount int

	// requestTimeout is the timeout duration of a request
	requestTimeout time.Duration

	log log15.Logger

	bus ari.Bus

	cancel context.CancelFunc
}

// New creates a new Client to the Asterisk ARI NATS proxy.
func New(ctx context.Context, opts ...OptionFunc) (ari.Client, error) {
	return newClient(ctx, opts...)
}

func newClient(ctx context.Context, opts ...OptionFunc) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &Client{
		cancel:            cancel,
		log:               log15.New(),
		prefix:            "ari.",
		uri:               "nats://localhost:4222",
		requestTimeout:    DefaultRequestTimeout,
		inputBufferLength: DefaultInputBufferLength,
	}
	c.log.SetHandler(log15.DiscardHandler())

	// Load environment-based configurations if such configurations are not explicitly set
	if os.Getenv("ARI_APPLICATION") != "" {
		c.appName = os.Getenv("ARI_APPLICATION")
	}
	if os.Getenv("NATS_URI") != "" {
		c.uri = os.Getenv("NATS_URI")
	}

	// Load explicit configurations
	for _, opt := range opts {
		opt(c)
	}

	// Make sure at least one of appName or dialog are set
	if c.appName == "" && c.dialog == "" {
		return nil, errors.New("at least one of Application and Dialog must be supplied")
	}

	// Connect to NATS, if we do not already have a connection
	if c.nc == nil {
		n, err := nats.Connect(c.uri)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to NATS")
		}
		c.nc, err = nats.NewEncodedConn(n, nats.JSON_ENCODER)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode NATS connection")
		}
		c.closeNATSOnClose = true
	}

	// Bind events from NATS to our bus
	// TODO: c.bus = stdbus.Start(ctx)
	go c.bindEvents(ctx)
	return c, nil
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
			c.nc = old.nc
			c.prefix = old.prefix
			c.log = old.log
			c.requestTimeout = old.requestTimeout

			// Make sure the old client does not close the NATS connection on us;
			if old.closeNATSOnClose {
				c.log.Warn("Disabling parent NATS connection closure; this will leak NATS connections")
				old.closeNATSOnClose = false
			}
		}
	}
}

// WithApplication configures the ARI Application to use for a Client
func WithApplication(name string) OptionFunc {
	return func(c *Client) {
		c.appName = name
	}
}

// WithAsterisk configures the ID of an Asterisk Node by which the Client should filter.
func WithAsterisk(id string) OptionFunc {
	return func(c *Client) {
		c.asterisk = id
	}
}

// WithDialog configures the Dialog ID to use for a Client
func WithDialog(id string) OptionFunc {
	return func(c *Client) {
		c.dialog = id
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

// WithPrefix configures the NATS Prefix to use om a Client
func WithPrefix(prefix string) OptionFunc {
	return func(c *Client) {
		c.prefix = prefix
	}
}

// ApplicationName returns the ARI application's name
func (c *Client) ApplicationName() string {
	return c.appName
}

// Close shuts down the client
func (c *Client) Close() {
	if c.cancel != nil {
		c.cancel()
	}

	if c.bus != nil {
		c.bus.Close()
	}

	if c.closeNATSOnClose && c.nc != nil {
		c.nc.Close()
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
	return nil
}

// Mailbox is the mailbox accessor
func (c *Client) Mailbox() ari.Mailbox {
	return &mailbox{c}
}

// Playback is the media playback accessor
func (c *Client) Playback() ari.Playback {
	return nil
}

// Sound is the sound accessor
func (c *Client) Sound() ari.Sound {
	return &sound{c}
}

// StoredRecording is the stored recording accessor
func (c *Client) StoredRecording() ari.StoredRecording {
	return nil
}

// TextMessage is the text message accessor
func (c *Client) TextMessage() ari.TextMessage {
	return nil
}

func (c *Client) commandRequest(req interface{}) error {
	var resp proxy.Response
	err := c.makeRequest(c.subject("command"), req, &resp)
	if err != nil {
		return err
	}
	return resp.Err()
}

func (c *Client) createRequest(req interface{}) (*proxy.Entity, error) {
	var resp proxy.Response
	err := c.makeRequest(c.subject("create"), req, &resp)
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

func (c *Client) getRequest(req interface{}) (*proxy.Entity, error) {
	var resp proxy.Response
	err := c.makeRequest(c.subject("get"), req, &resp)
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

func (c *Client) dataRequest(req interface{}) (*proxy.EntityData, error) {
	var resp proxy.Response
	err := c.makeRequest(c.subject("data"), req, &resp)
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

func (c *Client) listRequest(req interface{}) (*proxy.EntityList, error) {
	var resp proxy.Response
	err := c.makeRequest(c.subject("get"), req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, resp.Err()
	}
	if resp.EntityList == nil {
		return nil, ErrNil
	}
	return resp.EntityList, nil
}

func (c *Client) makeRequest(subject string, req interface{}, resp interface{}) error {
	return c.nc.Request(subject, req, resp, c.requestTimeout)
}

func (c *Client) subject(class string) (ret string) {
	ret = fmt.Sprintf("%s.%s", c.prefix, class)
	if c.appName != "" {
		ret += "." + c.appName
		if c.asterisk != "" {
			ret += "." + c.asterisk
		}
	}
	return
}

func (c *Client) eventsSubject() (subj string) {
	if c.appName != "" {
		subj = fmt.Sprintf("%sevent.%s", c.prefix, c.appName)

		if c.asterisk != "" {
			subj += "." + c.asterisk
		} else {
			subj += ".>"
		}
	}
	if c.dialog != "" {
		subj = fmt.Sprintf("%sdialogevent.%s", c.prefix, c.dialog)
	}
	return
}

func (c *Client) bindEvents(ctx context.Context) {
	subj := c.eventsSubject()
	if subj == "" {
		c.log.Error("cannot bind events without application or dialog")
		return
	}

	eChan := make(chan *ari.RawEvent, c.inputBufferLength)
	defer close(eChan)

	sub, err := c.nc.BindRecvChan(subj, eChan)
	if err != nil {
		c.log.Error("failed to bind NATS event receiver")
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case raw, ok := <-eChan:
			if !ok {
				c.log.Info("event channel closed")
				return
			}
			e, err := raw.ToEvent()
			if err != nil {
				c.log.Error("failed to convert raw event to ari.Event", "error", err)
				continue
			}
			c.bus.Send(e)
		}
	}
}
