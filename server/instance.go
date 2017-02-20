package ariproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari-proxy/server/dialog"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// Log is the internal logger for the ARI proxy.  It defaults to a no-op
// handler, but you may configure the handler at any time by calling
// `ariproxy.Log.SetHandler()`.  See the log15 documentation for details about
// the handler.
var Log log15.Logger

func init() {
	// Set up default (no-op) logger
	Log = log15.New()
	Log.SetHandler(log.DiscardHandler())
}

// Server describes the asterisk-facing ARI proxy server
type Server struct {
	// Application is the name of the ARI application of this server
	Application string

	// AsteriskID is the unique identifier for the Asterisk box
	// to which this server is connected.
	AsteriskID string

	// NATSPrefix is the string which should be prepended to all NATS subjects, sending and receiving.  It defaults to "ari.".
	NATSPrefix string

	// ari is the native Asterisk ARI client by which this proxy is directly connected
	ari *ari.Client

	// nats is the JSON-encoded NATS connection
	nats *nats.EncodedConn

	// Dialog is the dialog manager
	Dialog dialog.Manager

	readyCh chan struct{}

	// cancel is the context cancel function, by which all subtended subscriptions may be terminated
	cancel context.CancelFunc
}

// New returns a new Server
func New() *Server {
	return &Server{
		NATSPrefix: "ari.",
		readyCh:    make(chan struct{}),
		Dialog:     dialog.NewMemManager(),
	}
}

// Listen runs the given server, listening to ARI and NATS, as specified
func (s *Server) Listen(ctx context.Context, ariOpts native.Options, natsURI string) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// Connect to ARI
	s.ari, err = native.New(ariOpts)
	if err != nil {
		return errors.Wrap(err, "failed to connect to ARI")
	}
	defer s.ari.Close()

	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		return errors.Wrap(err, "failed to connect to NATS")
	}
	s.nats, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return errors.Wrap(err, "failed to encode NATS connection")
	}
	defer s.nats.Close()

	return s.listen(ctx)
}

// ListenOn runs the given server, listening on the provided ARI and NATS connections
func (s *Server) ListenOn(ctx context.Context, a *ari.Client, n *nats.EncodedConn) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.ari = a
	s.nats = n

	s.listen(ctx)
}

// Ready returns a channel which is closed when the Server is ready
func (s *Server) Ready() <-chan struct{} {
	if s.readyCh == nil {
		s.readyCh = make(chan struct{})
	}
	return s.readyCh
}

func (s *Server) listen(ctx context.Context) error {
	// First, get the Asterisk ID
	ret, err := s.ari.Asterisk.Info("")
	if err != nil {
		return errors.Wrap(err, "failed to get Asterisk ID")
	}

	s.AsteriskID = ret.SystemInfo.EntityID
	if s.AsteriskID == "" {
		return errors.New("empty Asterisk ID")
	}

	// Store the ARI application name for top-level access
	s.Application = s.ari.ApplicationName

	//
	// Listen on the initial NATS subjects
	//

	// ping handler
	pingSub, err := s.nats.Subscribe(fmt.Sprintf("%sping", s.NATSPrefix), s.pingHandler)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to pings")
	}
	defer pingSub.Unsubscribe()

	// get handlers
	allGet, err := s.nats.Subscribe(fmt.Sprintf("%sget", s.NATSPrefix), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-all subscription")
	}
	defer allGet.Unsubscribe()
	appGet, err := s.nats.Subscribe(fmt.Sprintf("%sget.%s", s.NATSPrefix, s.Application), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-app subscription")
	}
	defer appGet.Unsubscribe()
	idGet, err := s.nats.Subscribe(fmt.Sprintf("%sget.%s.%s", s.NATSPrefix, s.Application, s.AsteriskID), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-id subscription")
	}
	defer idGet.Unsubscribe()

	// command handlers
	allCommand, err := s.nats.Subscribe(fmt.Sprintf("%scommand", s.NATSPrefix), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-all subscription")
	}
	defer allCommand.Unsubscribe()
	appCommand, err := s.nats.Subscribe(fmt.Sprintf("%scommand.%s", s.NATSPrefix, s.Application), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-app subscription")
	}
	defer appCommand.Unsubscribe()
	idCommand, err := s.nats.Subscribe(fmt.Sprintf("%scommand.%s.%s", s.NATSPrefix, s.Application, s.AsteriskID), s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-id subscription")
	}
	defer idCommand.Unsubscribe()

	// create handlers
	allCreate, err := s.nats.QueueSubscribe(fmt.Sprintf("%screate", s.NATSPrefix), "ariproxy", s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-all subscription")
	}
	defer allCreate.Unsubscribe()
	appCreate, err := s.nats.QueueSubscribe(fmt.Sprintf("%screate.%s", s.NATSPrefix, s.Application), "ariproxy", s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-app subscription")
	}
	defer appCreate.Unsubscribe()
	idCreate, err := s.nats.QueueSubscribe(fmt.Sprintf("%screate.%s.%s", s.NATSPrefix, s.Application, s.AsteriskID), "ariproxy", s.requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-id subscription")
	}
	defer idCreate.Unsubscribe()

	// Run the periodic announcer
	go s.runAnnouncer(ctx)

	// Close the readyChannel to indicate that we are operational
	if s.readyCh != nil {
		close(s.readyCh)
	}

	// Wait for context closure to exit
	<-ctx.Done()
	return ctx.Err()
}

// runAnnouncer runs the periodic discovery announcer
func (s *Server) runAnnouncer(ctx context.Context) {
	ticker := time.NewTicker(proxy.AnnouncementInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.announce()
		}
	}
}

// announce publishes the presence of this server to the cluster
func (s *Server) announce() {
	s.nats.Publish(fmt.Sprintf("%sannounce", s.NATSPrefix), &proxy.Announcement{
		Asterisk:    s.AsteriskID,
		Application: s.Application,
	})
}

// pingHandler publishes the server's presence
func (s *Server) pingHandler(m *Msg) {
	s.announce()
}

func (s *Server) requestHandler(subject string, reply string, req *proxy.Request) {
	go s.dispatchRequest(req, reply)
}

func (s *Server) dispatchRequest(reply string, req *proxy.Request) {
	var f func(string, *proxy.Request)
	f := func(reply string, req *proxy.Request) {
		s.nats.Publish(reply, &proxy.ErrorResponse{Error: "Not implemented"})
	}

	if req.ApplicationList != nil {
		f = s.applicationList
	}
	if req.ApplicationData != nil {
		f = s.applicationData
	}
	if req.ApplicationSubscribe != nil {
		f = s.applicationSubscribe
	}
	if req.ApplicationUnsubscribe != nil {
		f = s.applicationUnsubscribe
	}

	return f(reply, req)
}

func (s *Server) sendError(reply string, err error) {
	s.nats.Publish(reply, proxy.NewErrorResponse(err))
}

func (s *Server) sendNotFound(reply string) {
	s.nats.Publish(reply, proxy.NewErrorResponse(ErrNotFound))
}

// Metadata returns the metadata for this server.  The dialog parameter is
// optional; set it to the empty string if one is not applicable or known.
func (s *Server) Metadata(dialog string) *proxy.Metadata {
	return &proxy.Metadata{
		Application: s.Application,
		Asterisk:    s.AsteriskID,
		Dialog:      dialog,
	}
}

func (s *Server) newInstance(id string, transport session.Transport) *Instance {
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

	i.log.Debug("Starting dialog instance")

	go func() {
		i.application()
		i.asterisk()
		i.modules()
		i.channel()
		i.storedRecording()
		i.liveRecording()
		i.bridge()
		i.device()
		i.playback()
		i.mailbox()
		i.sound()
		i.logging()
		i.config()

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

func (i *Instance) String() string {
	return fmt.Sprintf("Instance{%s}", i.Dialog.ID)
}
