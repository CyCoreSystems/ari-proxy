package server

import (
	"context"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/nats-io/nats"
)

func clientFactory(ctx context.Context) (ari.Client, error) {
	cl, err := client.New(ctx, client.WithApplication("asdf"), client.WithDialog("1234"))
	return cl, err
}

type srv struct {
	s *Server
}

func (s *srv) Start(ctx context.Context, t *testing.T, client ari.Client, nc *nats.EncodedConn, completeCh chan struct{}) {
	s.s = New()

	go func() {
		if err := s.s.ListenOn(ctx, client, nc); err != nil {
			if err != context.Canceled {
				t.Errorf("Failed to start server: %s", err)
			}
		}
		close(completeCh)
	}()
}

func (s *srv) Ready() <-chan struct{} {
	return s.s.Ready()
}

func (s *srv) Close() error {
	return nil
}
