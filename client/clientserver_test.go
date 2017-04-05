package client

import (
	"context"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/server"
	"github.com/nats-io/nats"
)

func clientFactory(ctx context.Context) (ari.Client, error) {
	cl, err := New(ctx, WithApplication("asdf"), WithDialog("1234"))
	return cl, err
}

type srv struct {
	s *server.Server
}

func (s *srv) Start(ctx context.Context, t *testing.T, client ari.Client, nc *nats.EncodedConn) {
	s.s = server.New()

	go func() {
		if err := s.s.ListenOn(ctx, client, nc); err != nil {
			if err != context.Canceled {
				t.Errorf("Failed to start server: %s", err)
			}
		}
	}()
}

func (s *srv) Ready() <-chan struct{} {
	return s.s.Ready()
}

func (s *srv) Close() error {
	return nil
}
