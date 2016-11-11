package client

import (
	"context"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/CyCoreSystems/ari/stdbus"
	"github.com/nats-io/nats"
)

// DefaultRequestTimeout is the default timeout for a NATS request
const DefaultRequestTimeout = 200 * time.Millisecond

// Options is the list options
type Options struct {
	// ReadOperationRetryCount is the amount of times to retry a read operation
	ReadOperationRetryCount int

	// RequestTimeout is the timeout duration of a request
	RequestTimeout time.Duration

	Parent context.Context
}

// New creates a new ari.Client connected to a gateway ARI server via NATS
func New(nc *nats.Conn, application string, d *session.Dialog, opts Options) (cl *ari.Client, err error) {
	if opts.RequestTimeout == 0 {
		opts.RequestTimeout = DefaultRequestTimeout
	}

	if opts.Parent == nil {
		opts.Parent = context.Background()
	}

	conn := &Conn{
		application: application,
		opts:        opts,
		conn:        nc,
		dialog:      d,
	}

	bus := stdbus.Start(opts.Parent)

	playback := natsPlayback{conn, bus}
	liveRecording := &natsLiveRecording{conn}
	storedRecording := &natsStoredRecording{conn}
	logging := &natsLogging{conn}
	modules := &natsModules{conn}
	config := &natsConfig{conn}

	cl = &ari.Client{
		Cleanup:     func() error { nc.Close(); return nil },
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

	return

}
