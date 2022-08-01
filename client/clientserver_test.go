package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari-proxy/v5/server"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/rid"
	"github.com/nats-io/nats.go"
)

type srv struct {
	s *server.Server
}

func (s *srv) Start(ctx context.Context, t *testing.T, mockClient ari.Client, nc *nats.EncodedConn, completeCh chan struct{}) (ari.Client, error) {
	s.s = server.New()

	// tests may run in parallel so we don't want two separate proxy servers to conflict.
	s.s.MBPrefix = rid.New("") + "."

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
	case <-time.After(500 * time.Millisecond):
		return nil, errors.New("Timeout waiting for server ready")
	}

	cl, err := New(ctx, WithTimeoutRetries(4), WithPrefix(s.s.MBPrefix), WithApplication("asdf"))
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
