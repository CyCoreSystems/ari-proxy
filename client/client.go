package client

import (
	"context"
	"fmt"
	"os"
	"time"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// DefaultInputBufferLength is the default size of the event buffer for events coming in from NATS
const DefaultInputBufferLength = 100

// Client is the ari-proxy client.
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
	cl, err := newClient(ctx, opts...)
	return cl, err
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

	// Load explicit configurations
	for _, opt := range opts {
		opt(c)
	}

	// Load environment-based configurations if such configurations are not explicitly set
	if c.appName == "" && os.Getenv("ARI_APPLICATION") != "" {
		c.appName = os.Getenv("ARI_APPLICATION")
	}
	if c.uri == "" && os.Getenv("NATS_URI") != "" {
		c.uri = os.Getenv("NATS_URI")
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

func (p *Client) ApplicationName() string {
	return p.appName
}

func (p *Client) Close() {
	if p.cancel != nil {
		p.cancel()
	}

	if p.bus != nil {
		p.bus.Close()
	}

	if p.closeNATSOnClose && p.nc != nil {
		p.nc.Close()
	}
}

func (p *Client) Application() ari.Application {
	return &application{
		c: p,
	}
}

func (p *Client) Asterisk() ari.Asterisk {
	return &asterisk{
		c: p,
	}
}

func (p *Client) Bridge() ari.Bridge {
	return &bridge{
		c: p,
	}
}

func (p *Client) Bus() ari.Bus {
	return p.bus
}

func (p *Client) Channel() ari.Channel {
	return &channel{p}
}

func (p *Client) DeviceState() ari.DeviceState {
	return nil
}

func (p *Client) Endpoint() ari.Endpoint {
	return nil
}

func (p *Client) LiveRecording() ari.LiveRecording {
	return &liveRecording{p}
}

func (p *Client) Mailbox() ari.Mailbox {
	return nil
}

func (p *Client) Playback() ari.Playback {
	return &playback{p}
}

func (p *Client) Sound() ari.Sound {
	return nil
}

func (p *Client) StoredRecording() ari.StoredRecording {
	return nil
}

func (p *Client) TextMessage() ari.TextMessage {
	return nil
}

func (p *Client) subject(class string) (ret string) {
	ret = fmt.Sprintf("%s.%s", p.prefix, class)
	if p.appName != "" {
		ret += "." + p.appName
		if p.asterisk != "" {
			ret += "." + p.asterisk
		}
	}
	return
}

func (p *Client) eventsSubject() (subj string) {
	if p.appName != "" {
		subj = fmt.Sprintf("%sevent.%s", p.prefix, p.appName)

		if p.asterisk != "" {
			subj += "." + p.asterisk
		} else {
			subj += ".>"
		}
	}
	if p.dialog != "" {
		subj = fmt.Sprintf("%sdialogevent.%s", p.prefix, p.dialog)
	}
	return
}

func (p *Client) bindEvents(ctx context.Context) {
	subj := p.eventsSubject()
	if subj == "" {
		p.log.Error("cannot bind events without application or dialog")
		return
	}

	eChan := make(chan *ari.RawEvent, p.inputBufferLength)
	defer close(eChan)

	sub, err := p.nc.BindRecvChan(subj, eChan)
	if err != nil {
		p.log.Error("failed to bind NATS event receiver")
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case raw, ok := <-eChan:
			if !ok {
				p.log.Info("event channel closed")
				return
			}
			e, err := raw.ToEvent()
			if err != nil {
				p.log.Error("failed to convert raw event to ari.Event", "error", err)
				continue
			}
			p.bus.Send(e)
		}
	}
}
