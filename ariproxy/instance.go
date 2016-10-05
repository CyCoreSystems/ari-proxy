package ariproxy

import (
	"sync"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
	"golang.org/x/net/context"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// An Instance is the server-side component of a dialog
type Instance struct {
	Dialog *session.Dialog

	server *Server

	dispatcher     map[string]Handler2
	dispatcherLock sync.RWMutex

	readyCh chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	upstream *ari.Client
	conn     *nats.Conn
	log      log15.Logger
}

func (srv *Server) newInstance(id string, transport session.Transport) *Instance {
	return &Instance{
		Dialog:     session.NewDialog(id, transport),
		readyCh:    make(chan struct{}),
		server:     srv,
		upstream:   srv.upstream,
		conn:       srv.conn,
		log:        srv.log.New("dialog", id),
		dispatcher: make(map[string]Handler2),
	}
}

// Start runs the server side instance
func (i *Instance) Start(ctx context.Context) {
	i.ctx, i.cancel = context.WithCancel(ctx)

	go func() {
		i.application()
		i.asterisk()
		i.modules()

		/*
			i.bridge()
			i.channel()
			i.device()
			i.playback()
			i.mailbox()
			i.sound()
			i.liveRecording()
			i.storedRecording()
			i.logging()
			i.config()
		*/

		// do commands last, since that is the one that will be dispatching
		i.commands()

		close(i.readyCh)

		<-i.ctx.Done()
	}()

	<-i.readyCh
}

// Stop stops the instance
func (i *Instance) Stop() {
	if i == nil {
		return
	}
	i.cancel()
}
