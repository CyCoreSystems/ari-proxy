package client

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/stdbus"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// Options is the list options
type Options struct {
	// Application provides an ARI application to which to bind.  At least one of Application or Dialog must be specified.  This option may also be supplied by the `ARI_APPLICATION` environment variable.
	Application string

	// Dialog provides a dialog ID to which to bind.  At least one of Application or Dialog must be specified.
	Dialog string

	// NATS provides the nats.EncodedConn over which messages will be transceived.  One of NATS or NATSURI must be specified.
	NATS *nats.EncodedConn

	// NATSURI provies the URI to which a NATS connection should be established. One of NATS or NATSURI must be specified. This option may also be supplied by the `NATS_URI` environment variable.
	NATSURI string

	// NATSPrefix is the prefix to use on all NATS subjects.  It defaults to "ari.".
	NATSPrefix string

	// closeNATSOnClose indicates that the NATS connection should be closed when the ari.Client is closed
	closeNATSOnClose bool

	// ReadOperationRetryCount is the amount of times to retry a read operation
	ReadOperationRetryCount int

	// RequestTimeout is the timeout duration of a request
	RequestTimeout time.Duration
}

type proxyClient struct {
	// application provides an ARI application to which to bind
	application string

	// dialog provides a dialog ID to which to bind
	dialog string

	// nats provides the nats.EncodedConn over which messages will be transceived
	nats *nats.EncodedConn

	application string
	dialog      string

	nats *nats.EncodedConn
}

// New creates a new proxy-based ari.Client bound to the given application
func New(ctx context.Context, nc *nats.EncodedConn, opts Options) (cl *ari.Client, err error) {
	if opts.Application == "" {
		if os.Getenv("ARI_APPLICATION") == "" {
			return nil, errors.New("at least one of Application and Dialog must be supplied")
		}
		opts.Application = os.Getenv("ARI_APPLICATION")
	}

	var nc *nats.EncodedConn
	if opts.NATS == nil {
		if opts.NATSURI == "" {
			if os.Getenv("NATS_URI") == "" {
				return nil, errors.New("one of NATS or NATSURI must be supplied")
			}
			opts.NATSURI = os.Getenv("NATS_URI")
		}

		n, err := nats.Connect(opts.NATSURI)
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to NATS")
		}
		nc, err = nats.NewEncodedConn(n, nats.JSON_ENCODER)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode NATS connection")
		}
		opts.closeNATSOnClose = true
	} else {
		nc = opts.NATS
	}

	if opts.NATSPrefix == "" {
		opts.NATSPrefix = "ari."
	}

	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = DefaultRequestTimeout
	}

	bus := stdbus.Start(ctx)

	playback := natsPlayback{conn, bus}
	liveRecording := &natsLiveRecording{conn}
	storedRecording := &natsStoredRecording{conn}
	logging := &natsLogging{conn}
	modules := &natsModules{conn}
	config := &natsConfig{conn}

	cl = &ari.Client{
		Asterisk:    &natsAsterisk{conn, logging, modules, config},
		Application: &natsApplication{conn},
		Bridge:      &natsBridge{conn, bus, &playback, liveRecording},
		Channel:     &natsChannel{conn, bus, &playback, liveRecording},
		DeviceState: &natsDeviceState{conn},
		Mailbox:     &natsMailbox{conn},
		Sound:       &natsSound{conn},
		Playback:    &playback,
		Recording: &ari.Recording{
			Live:   liveRecording,
			Stored: storedRecording,
		},
		Bus:             bus,
		ApplicationName: application,
	}

	if opts.closeNATSOnClose {
		cl.Cleanup = func() {
			nc.Close()
		}
	}

	return

}
