package integration

import (
	"context"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
)

func natsConnect() (*nats.EncodedConn, error) {
	c, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to NATS")
	}
	nc, err := nats.NewEncodedConn(c, nats.JSON_ENCODER)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode NATS connection")
	}
	return nc, err
}

// these hooks are designed to allow both client and server packages
// to call this integration package without recursive imports

// ClientFactory creates a client
type ClientFactory func(ctx context.Context) (ari.Client, error)

// Server represents a generalized ari-proxy server
type Server interface {
	Start(ctx context.Context, t *testing.T, client ari.Client, nc *nats.EncodedConn)
	Ready() <-chan struct{}
	Close() error
}

func runTest(desc string, t *testing.T, s Server, clientFactory ClientFactory, fn func(m *mock, cl ari.Client)) {
	t.Run(desc, func(t *testing.T) {
		// setup mocking
		m := standardMock()

		// setup ari-proxy server
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		nc, err := natsConnect()
		if err != nil {
			t.Skipf("Error connecting to nats: %s", err)
		}

		s.Start(ctx, t, m.Client, nc)

		select {
		case <-s.Ready():
		case <-time.After(2 * time.Second):
			t.Errorf("Timeout waiting for server to be ready")
		}

		// setup client
		cl, err := clientFactory(ctx)
		if err != nil {
			t.Errorf("Failed to build ari-proxy client: %s", err)
			return
		}
		defer cl.Close()

		fn(m, cl)
	})
}
