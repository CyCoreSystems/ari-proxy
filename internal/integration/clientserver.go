package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"sync"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats.go"
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

// Server represents a generalized ari-proxy server
type Server interface {
	Start(ctx context.Context, t *testing.T, client ari.Client, nc *nats.EncodedConn, completeCh chan struct{}) (ari.Client, error)
	Ready() <-chan struct{}
	Close() error
}

// TestHandler is the interface for test execution
type testHandler func(t *testing.T, m *mock, cl ari.Client)

func runTest(desc string, t *testing.T, s Server, fn testHandler) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if err := recover(); err != nil {
			t.Errorf("PANIC")
		}
	}()

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

		completeCh := make(chan struct{})

		cl, err := s.Start(ctx, t, m.Client, nc, completeCh)
		if err != nil {
			t.Errorf("Failed to start client/server: %s", err)
			return
		}
		defer cl.Close()

		var once sync.Once

		m.Shutdown = func() {
			once.Do(func() {
				s.Close()

				cancel()
				<-completeCh
			})
		}

		fn(t, m, cl)

		m.Shutdown()

		timeout, ok := TimeoutCount(cl)
		if !ok {
			t.Errorf("Failed to get timeout count from ari-proxy client")
		} else {
			if timeout > 0 {
				fmt.Fprintf(os.Stderr, "Timeouts: %d\n", timeout)
			}
		}

	})
}

type timeoutCounter interface {
	TimeoutCount() int64
}

// TimeoutCount gets the timeout count from the ari client, if available.
func TimeoutCount(c ari.Client) (int64, bool) {
	cl, ok := c.(timeoutCounter)
	if !ok {
		return 0, false
	}
	return cl.TimeoutCount(), true
}
