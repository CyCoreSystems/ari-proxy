package server

import (
	"context"
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari-proxy/server/dialog"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

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
	ari ari.Client

	// nats is the JSON-encoded NATS connection
	nats *nats.EncodedConn

	// Dialog is the dialog manager
	Dialog dialog.Manager

	readyCh chan struct{}

	// cancel is the context cancel function, by which all subtended subscriptions may be terminated
	cancel context.CancelFunc

	// Log is the log15.Logger for the service.  You may replace or call SetHandler() on this at any time to change the logging of the service.
	Log log15.Logger
}

// New returns a new Server
func New() *Server {
	log := log15.New()
	log.SetHandler(log15.DiscardHandler())

	return &Server{
		NATSPrefix: "ari.",
		readyCh:    make(chan struct{}),
		Dialog:     dialog.NewMemManager(),
		Log:        log,
	}
}

// Listen runs the given server, listening to ARI and NATS, as specified
func (s *Server) Listen(ctx context.Context, ariOpts native.Options, natsURI string) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// Connect to ARI
	s.ari, err = native.Connect(ctx, ariOpts)
	if err != nil {
		return errors.Wrap(err, "failed to connect to ARI")
	}
	defer s.ari.Close()

	// Connect to NATS
	nc, err := nats.Connect(natsURI)
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
func (s *Server) ListenOn(ctx context.Context, a ari.Client, n *nats.EncodedConn) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.ari = a
	s.nats = n

	return s.listen(ctx)
}

// Ready returns a channel which is closed when the Server is ready
func (s *Server) Ready() <-chan struct{} {
	if s.readyCh == nil {
		s.readyCh = make(chan struct{})
	}
	return s.readyCh
}

// nolint: gocyclo
func (s *Server) listen(ctx context.Context) error {

	var wg closeGroup
	defer func() {
		select {
		case <-wg.Done():
		case <-time.After(500 * time.Millisecond):
			panic("timeout waiting for shutdown of sub components")
		}
	}()

	// First, get the Asterisk ID

	ret, err := s.ari.Asterisk().Info("")
	if err != nil {
		return errors.Wrap(err, "failed to get Asterisk ID")
	}

	s.AsteriskID = ret.SystemInfo.EntityID
	if s.AsteriskID == "" {
		return errors.New("empty Asterisk ID")
	}

	// Store the ARI application name for top-level access
	s.Application = s.ari.ApplicationName()

	//
	// Listen on the initial NATS subjects
	//

	// ping handler
	pingSub, err := s.nats.Subscribe(fmt.Sprintf("%sping", s.NATSPrefix), s.pingHandler)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to pings")
	}
	defer wg.Add(pingSub.Unsubscribe)

	// get a contextualized request handler
	requestHandler := s.newRequestHandler(ctx)

	// get handlers
	allGet, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "get", "", ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-all subscription")
	}
	defer wg.Add(allGet.Unsubscribe)()

	appGet, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "get", s.Application, ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-app subscription")
	}
	defer wg.Add(appGet.Unsubscribe)()
	idGet, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "get", s.Application, s.AsteriskID), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create get-id subscription")
	}
	defer wg.Add(idGet.Unsubscribe)()

	// data handlers
	allData, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "data", "", ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create data-all subscription")
	}
	defer wg.Add(allData.Unsubscribe)()
	appData, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "data", s.Application, ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create data-app subscription")
	}
	defer wg.Add(appData.Unsubscribe)()
	idData, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "data", s.Application, s.AsteriskID), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create data-id subscription")
	}
	defer wg.Add(idData.Unsubscribe)()

	// command handlers
	allCommand, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "command", "", ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-all subscription")
	}
	defer wg.Add(allCommand.Unsubscribe)()
	appCommand, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "command", s.Application, ""), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-app subscription")
	}
	defer wg.Add(appCommand.Unsubscribe)()
	idCommand, err := s.nats.Subscribe(proxy.Subject(s.NATSPrefix, "command", s.Application, s.AsteriskID), requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create command-id subscription")
	}
	defer wg.Add(idCommand.Unsubscribe)()

	// create handlers
	allCreate, err := s.nats.QueueSubscribe(proxy.Subject(s.NATSPrefix, "create", "", ""), "ariproxy", requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-all subscription")
	}
	defer wg.Add(allCreate.Unsubscribe)()
	appCreate, err := s.nats.QueueSubscribe(proxy.Subject(s.NATSPrefix, "create", s.Application, ""), "ariproxy", requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-app subscription")
	}
	defer wg.Add(appCreate.Unsubscribe)()
	idCreate, err := s.nats.QueueSubscribe(proxy.Subject(s.NATSPrefix, "create", s.Application, s.AsteriskID), "ariproxy", requestHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create create-id subscription")
	}
	defer wg.Add(idCreate.Unsubscribe)()

	// Run the periodic announcer
	go s.runAnnouncer(ctx)

	// Run the event handler
	go s.runEventHandler(ctx)

	// TODO: run the dialog cleanup routine (remove bindings for entities which no longer exist)
	//go s.runDialogCleaner(ctx)

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

// runEventHandler processes events which are received from ARI
func (s *Server) runEventHandler(ctx context.Context) {
	sub := s.ari.Bus().Subscribe(ari.Events.All)
	defer sub.Cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case e := <-sub.Events():
			raw := ari.EventToRaw(e)
			if raw == nil {
				s.Log.Error("Failed to encode event to RawEvent")
				continue
			}

			// Add metadata
			raw.Header.Set("application", s.Application)
			raw.Header.Set("asterisk", s.AsteriskID)

			// Publish event to canonical destination
			s.nats.Publish(fmt.Sprintf("%sevent.%s.%s", s.NATSPrefix, s.Application, s.AsteriskID), raw)

			// Publish event to any associated dialogs
			for _, d := range s.dialogsForEvent(e) {
				dRaw := raw
				dRaw.Header.Set("dialog", d)
				s.nats.Publish(fmt.Sprintf("%sdialogevent.%s", s.NATSPrefix, d), dRaw)
			}
		}
	}
}

// pingHandler publishes the server's presence
func (s *Server) pingHandler(m *nats.Msg) {
	s.announce()
}

// newRequestHandler returns a context-wrapped nats.Handler to handle requests
func (s *Server) newRequestHandler(ctx context.Context) func(subject string, reply string, req *proxy.Request) {
	return func(subject string, reply string, req *proxy.Request) {
		go s.dispatchRequest(ctx, reply, req)
	}
}

// nolint: gocyclo
func (s *Server) dispatchRequest(ctx context.Context, reply string, req *proxy.Request) {
	f := func(ctx context.Context, reply string, req *proxy.Request) {
		s.sendError(reply, errors.New("Not implemented"))
	}

	if req.Metadata == nil {
		req.Metadata = &proxy.Metadata{}
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

	if req.AsteriskInfo != nil {
		f = s.asteriskInfo
	}
	if req.AsteriskReloadModule != nil {
		f = s.asteriskReloadModule
	}
	if req.AsteriskVariables != nil {
		if req.AsteriskVariables.Get != nil {
			f = s.asteriskVariableGet
		}

		if req.AsteriskVariables.Set != nil {
			f = s.asteriskVariableSet
		}
	}

	if req.BridgeAddChannel != nil {
		f = s.bridgeAddChannel
	}
	if req.BridgeCreate != nil {
		f = s.bridgeCreate
	}
	if req.BridgeData != nil {
		f = s.bridgeData
	}
	if req.BridgeList != nil {
		f = s.bridgeList
	}
	if req.BridgePlay != nil {
		f = s.bridgePlay
	}
	if req.BridgeDelete != nil {
		f = s.bridgeDelete
	}
	if req.BridgeRecord != nil {
		f = s.bridgeRecord
	}
	if req.BridgeRemoveChannel != nil {
		f = s.bridgeRemoveChannel
	}
	if req.BridgeSubscribe != nil {
		f = s.bridgeSubscribe
	}

	if req.ChannelAnswer != nil {
		f = s.channelAnswer
	}
	if req.ChannelBusy != nil {
		f = s.channelBusy
	}
	if req.ChannelCongestion != nil {
		f = s.channelCongestion
	}
	if req.ChannelCreate != nil {
		f = s.channelCreate
	}
	if req.ChannelData != nil {
		f = s.channelData
	}
	if req.ChannelContinue != nil {
		f = s.channelContinue
	}
	if req.ChannelDial != nil {
		f = s.channelDial
	}
	if req.ChannelHangup != nil {
		f = s.channelHangup
	}
	if req.ChannelHold != nil {
		f = s.channelHold
	}
	if req.ChannelList != nil {
		f = s.channelList
	}
	if req.ChannelMOH != nil {
		f = s.channelMOH
	}
	if req.ChannelMute != nil {
		f = s.channelMute
	}
	if req.ChannelOriginate != nil {
		f = s.channelOriginate
	}
	if req.ChannelPlay != nil {
		f = s.channelPlay
	}
	if req.ChannelRecord != nil {
		f = s.channelRecord
	}
	if req.ChannelRing != nil {
		f = s.channelRing
	}
	if req.ChannelSendDTMF != nil {
		f = s.channelSendDTMF
	}
	if req.ChannelSilence != nil {
		f = s.channelSilence
	}
	if req.ChannelSnoop != nil {
		f = s.channelSnoop
	}
	if req.ChannelStopHold != nil {
		f = s.channelStopHold
	}
	if req.ChannelStopMOH != nil {
		f = s.channelStopMOH
	}
	if req.ChannelStopRing != nil {
		f = s.channelStopRing
	}
	if req.ChannelStopSilence != nil {
		f = s.channelStopSilence
	}
	if req.ChannelSubscribe != nil {
		f = s.channelSubscribe
	}
	if req.ChannelUnmute != nil {
		f = s.channelUnmute
	}
	if req.ChannelVariables != nil {
		if req.ChannelVariables.Get != nil {
			f = s.channelVariableGet
		}

		if req.ChannelVariables.Set != nil {
			f = s.channelVariableSet
		}
	}

	if req.DeviceStateData != nil {
		f = s.deviceStateData
	}
	if req.DeviceStateDelete != nil {
		f = s.deviceStateDelete
	}
	if req.DeviceStateList != nil {
		f = s.deviceStateList
	}
	if req.DeviceStateUpdate != nil {
		f = s.deviceStateUpdate
	}

	if req.EndpointData != nil {
		f = s.endpointData
	}
	if req.EndpointList != nil {
		f = s.endpointList
	}
	if req.EndpointListByTech != nil {
		f = s.endpointListByTech
	}

	if req.MailboxData != nil {
		f = s.mailboxData
	}
	if req.MailboxDelete != nil {
		f = s.mailboxDelete
	}
	if req.MailboxList != nil {
		f = s.mailboxList
	}
	if req.MailboxUpdate != nil {
		f = s.mailboxUpdate
	}

	if req.PlaybackControl != nil {
		f = s.playbackControl
	}
	if req.PlaybackData != nil {
		f = s.playbackData
	}
	if req.PlaybackStop != nil {
		f = s.playbackStop
	}
	if req.PlaybackSubscribe != nil {
		f = s.playbackSubscribe
	}

	if req.RecordingStoredCopy != nil {
		f = s.recordingStoredCopy
	}
	if req.RecordingStoredData != nil {
		f = s.recordingStoredData
	}
	if req.RecordingStoredDelete != nil {
		f = s.recordingStoredDelete
	}
	if req.RecordingStoredList != nil {
		f = s.recordingStoredList
	}

	if req.RecordingLiveData != nil {
		f = s.recordingLiveData
	}
	if req.RecordingLiveDelete != nil {
		f = s.recordingLiveDelete
	}
	if req.RecordingLiveMute != nil {
		f = s.recordingLiveMute
	}
	if req.RecordingLivePause != nil {
		f = s.recordingLivePause
	}
	if req.RecordingLiveResume != nil {
		f = s.recordingLiveResume
	}
	if req.RecordingLiveScrap != nil {
		f = s.recordingLiveScrap
	}
	if req.RecordingLiveStop != nil {
		f = s.recordingLiveStop
	}
	if req.RecordingLiveUnmute != nil {
		f = s.recordingLiveUnmute
	}

	if req.SoundData != nil {
		f = s.soundData
	}
	if req.SoundList != nil {
		f = s.soundList
	}

	if req.AsteriskConfig != nil {
		if req.AsteriskConfig.Data != nil {
			f = s.asteriskConfigData
		}

		if req.AsteriskConfig.Delete != nil {
			f = s.asteriskConfigDelete
		}

		if req.AsteriskConfig.Update != nil {
			f = s.asteriskConfigUpdate
		}
	}

	if req.AsteriskLogging != nil {
		if req.AsteriskLogging.List != nil {
			f = s.asteriskLoggingList
		}

		if req.AsteriskLogging.Create != nil {
			f = s.asteriskLoggingCreate
		}

		if req.AsteriskLogging.Data != nil {
			f = s.asteriskLoggingData
		}

		if req.AsteriskLogging.Rotate != nil {
			f = s.asteriskLoggingRotate
		}

		if req.AsteriskLogging.Delete != nil {
			f = s.asteriskLoggingDelete
		}
	}

	if req.AsteriskModules != nil {
		if req.AsteriskModules.List != nil {
			f = s.asteriskModuleList
		}

		if req.AsteriskModules.Data != nil {
			f = s.asteriskModuleData
		}

		if req.AsteriskModules.Load != nil {
			f = s.asteriskModuleLoad
		}

		if req.AsteriskModules.Unload != nil {
			f = s.asteriskModuleUnload
		}

		if req.AsteriskModules.Reload != nil {
			f = s.asteriskModuleReload
		}
	}

	f(ctx, reply, req)
}

func (s *Server) sendError(reply string, err error) {
	s.nats.Publish(reply, proxy.NewErrorResponse(err))
}

func (s *Server) sendNotFound(reply string) {
	s.nats.Publish(reply, proxy.NewErrorResponse(proxy.ErrNotFound))
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

/*
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
*/
