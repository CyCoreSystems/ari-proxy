package server

import (
	"context"
	"errors"
	"testing"

	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

type srv struct {
	s *Server
}

func (s *srv) Start(ctx context.Context, t *testing.T, mockClient ari.Client, nc *nats.EncodedConn, completeCh chan struct{}) (ari.Client, error) {
	s.s = New()
	// tests may run in parallel so we don't want two separate proxy servers to conflict.
	s.s.NATSPrefix = uuid.NewV1().String() + "."
	s.s.Application = "asdf"

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
	case <-time.After(600 * time.Millisecond):
		return nil, errors.New("Timeout waiting for server ready")
	}

	cl, err := client.New(ctx, client.WithTimeoutRetries(4), client.WithPrefix(s.s.NATSPrefix), client.WithApplication("asdf"))
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
