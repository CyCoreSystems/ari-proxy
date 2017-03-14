package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari/stdbus"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// DefaultInputBufferLength is the default size of the event buffer for events coming in from NATS
const DefaultInputBufferLength = 100

type proxyClient struct {

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
	ctx, cancel := context.WithCancel(ctx)

	c := &proxyClient{
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
	if c.appName == "" && os.Getenv("ARI_APPLICATION") {
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
	c.bus = stdbus.Start(ctx)
	go c.bindEvents(ctx)
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
		old, ok := cl.(*proxyClient)
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
		c.NATS = nc
	}
}

// WithNATSURI configures the NATS URI for a Client
func WithNATSURI(uri string) OptionFunc {
	return func(c *Client) {
		c.NATSURI = nc
	}
}

// WithPrefix configures the NATS Prefix to use om a Client
func WithPrefix(prefix string) OptionFunc {
	return func(c *Client) {
		c.NATSPrefix = prefix
	}
}

func (p *proxyClient) ApplicationName() string {
	return p.appName
}

func (p *proxyClient) Close() {
	if p.cancel != nil {
		p.cancel()
	}

	if p.bus != nil {
		p.bus.Close()
	}

	if opts.closeNATSOnClose && p.nc != nil {
		p.nc.Close()
	}
}

func (p *proxyClient) Application() ari.Application {
	return &Application{p}
}

func (p *proxyClient) Asterisk() ari.Asterisk {
	return &Asterisk{p}
}

func (p *proxyClient) Bridge() ari.Bridge {
	return &Bridge{p}
}

func (p *proxyClient) Bus() ari.Bus {
	return p.bus
}

func (p *proxyClient) Channel() ari.Channel {
	return &Channel{p}
}

func (p *proxyClient) DeviceState() ari.DeviceState {
	return &DeviceState{p}
}

func (p *proxyClient) Endpoint() ari.Endpoint {
	return &Endpoint{p}
}

func (p *proxyClient) LiveRecording() ari.LiveRecording {
	return &LiveRecording{p}
}

func (p *proxyClient) Mailbox() ari.Mailbox {
	return &Mailbox{p}
}

func (p *proxyClient) Playback() ari.Playback {
	return &Playback{p}
}

func (p *proxyClient) Sound() ari.Sound {
	return &Sound{p}
}

func (p *proxyClient) StoredRecording() ari.StoredRecording {
	return &StoredRecording{p}
}

func (p *proxyClient) TextMessage() ari.TextMessage {
	return &TextMessage{p}
}

func (p *proxyClient) commandRequest(req interface{}) error {
	var resp proxy.Response
	err := p.makeRequest(subject("command"), req, &resp)
	if err != nil {
		return err
	}
	return resp.Err()
}

func (p *proxyClient) createRequest(req interface{}) (*proxy.Entity, error) {
	var resp proxy.Response
	err := p.makeRequest(subject, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, err
	}
	if resp.Entity == nil {
		return nil, errors.New("no entity returned")
	}
	return resp.Entity, nil
}

func (p *proxyClient) getRequest(req interface{}) (*proxy.Entity, error) {
	var resp proxy.Response
	err := p.makeRequest(subject, req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Err() != nil {
		return nil, err
	}
	if resp.Entity == nil {
		return nil, errors.New("no entity returned")
	}
	return resp.Entity, nil

}

func (p *proxyClient) dataRequest(req interface{}, resp interface{}) error {

}

func (p *proxyClient) listRequest(req interface{}) (*proxy.EntityList, error) {

}

func (p *proxyClient) makeRequest(subject string, req interface{}, resp interface{}) error {
	p.nc.Request(subject)
}

func (p *proxyClient) subject(class string) (ret string) {
	ret = fmt.Sprintf("%s.%s", p.prefix, class)
	if p.appName != "" {
		ret += "." + p.appName
		if p.asterisk != "" {
			ret += "." + p.asterisk
		}
	}
	return
}

func (p *proxyClient) eventsSubject() (subj string) {
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

func (p *proxyClient) bindEvents(ctx context.Context) {
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
