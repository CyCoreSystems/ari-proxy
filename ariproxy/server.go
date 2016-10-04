package ariproxy

import (
	"errors"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Server is the nats gateway server
type Server struct {
	ID          string // server identifier
	Application string // name of the asterisk application this gateway is serving

	cache instanceCache

	readyCh chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	upstream *ari.Client
	conn     *nats.Conn
	log      log15.Logger
}

// NewServer creates a new nats gw server
func NewServer(client *ari.Client, application string, opts *Options) (srv *Server, err error) {

	id := uuid.NewV1().String() //TODO: allow users to specify server, load from hostname, etc?

	if client == nil {
		err = errors.New("No client provided")
		return
	}

	if opts == nil {
		opts = &Options{}
	}

	if opts.Logger == nil {
		opts.Logger = log15.New()
	}

	if opts.Parent == nil {
		opts.Parent = context.Background()
	}

	if opts.URL == "" {
		opts.URL = nats.DefaultURL
	}

	srv = &Server{
		ID:          id,
		Application: application,
		readyCh:     make(chan struct{}),
		log:         opts.Logger,
		upstream:    client,
	}
	defer func() {
		if err != nil {
			srv = nil // don't return and garbage collect srv on error
		}
	}()

	srv.conn, err = nats.Connect(opts.URL)
	if err != nil {
		return
	}

	srv.ctx, srv.cancel = context.WithCancel(opts.Parent)

	return
}

// Start starts the service and listens for nats requests and delegates them to the upstream ARI client
func (srv *Server) Start() {

	go func() {
		defer srv.conn.Close()

		srv.events()

		close(srv.readyCh)

		<-srv.ctx.Done()
	}()

	<-srv.readyCh
}

// Close closes the gateway server
func (srv *Server) Close() {
	if srv == nil {
		return
	}
	srv.cancel()
}
