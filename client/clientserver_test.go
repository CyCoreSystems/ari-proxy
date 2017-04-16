package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/server"
	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

type srv struct {
	s *server.Server
}

func (s *srv) Start(ctx context.Context, t *testing.T, mockClient ari.Client, nc *nats.EncodedConn, completeCh chan struct{}) (ari.Client, error) {
	s.s = server.New()

	// tests may run in parallel so we don't want two separate proxy servers to conflict.
	s.s.NATSPrefix = uuid.NewV1().String() + "."

	go func() {
		if err := s.s.ListenOn(ctx, mockClient, nc); err != nil {
			if err != context.Canceled {
				t.Errorf("Failed to start server: %s", err)
			}
		}
		close(completeCh)
	}()

	select {
	case <-s.s.Ready():
	case <-time.After(200 * time.Millisecond):
		return nil, errors.New("Timeout waiting for server ready")
	}

	cl, err := New(ctx, WithTimeoutRetry(4), WithPrefix(s.s.NATSPrefix), WithApplication("asdf"), WithDialog("1234"))
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (s *srv) Ready() <-chan struct{} {
	return s.s.Ready()
}

func (s *srv) Close() error {
	return nil
}
