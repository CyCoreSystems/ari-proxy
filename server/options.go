package server

import (
	"context"

	"github.com/inconshreveable/log15"
)

// Options are the group of options for the ari-proxy server
type Options struct {
	//TODO: include nats/rabbitmq options

	URL string

	Logger log15.Logger
	Parent context.Context
}
