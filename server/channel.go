package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) channelAnswer(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelAnswer.ID
	err := s.ari.Channel().Answer(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelBusy(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel().Busy(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelCongestion(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel().Congestion(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelCreate(ctx context.Context, reply string, req *proxy.Request) {
	create := req.ChannelCreate.ChannelCreateRequest

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", create.ChannelID)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", create.OtherChannelID)
	}

	handle, err := s.ari.Channel().Create(create)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", handle.ID())
	}

	s.nats.Publish(reply, handle.ID())
}

func (s *Server) channelData(ctx context.Context, reply string, req *proxy.Request) {
	d, err := s.ari.Channel().Data(req.ChannelData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &d)
}

func (s *Server) channelContinue(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Continue(req.ChannelContinue.ID, req.ChannelContinue.Context, req.ChannelContinue.Extension, req.ChannelContinue.Priority))
}

func (s *Server) channelDial(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Dial(req.ChannelDial.ID, req.ChannelDial.Caller, req.ChannelDial.Timeout))
}

func (s *Server) channelHangup(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Hangup(req.ChannelHangup.ID, req.ChannelHangup.Reason))
}

func (s *Server) channelHold(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Hold(req.ChannelHold.ID))
}

func (s *Server) channelList(ctx context.Context, reply string, req *proxy.Request) {
	cx, err := s.ari.Channel().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var channels []string
	for _, channel := range cx {
		channels = append(channels, channel.ID())
	}

	s.nats.Publish(reply, &channels)
}

func (s *Server) channelMOH(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelMOH.ID
	music := req.ChannelMOH.Music
	err := s.ari.Channel().MOH(id, music)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelMute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Mute(req.ChannelMute.ID, req.ChannelMute.Direction))
}

func (s *Server) channelOriginate(ctx context.Context, reply string, req *proxy.Request) {
	orig := req.ChannelOriginate.OriginateRequest

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.ChannelID)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.OtherChannelID)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.Originator)
	}

	handle, err := s.ari.Channel().Originate(orig)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", handle.ID())
	}

	s.nats.Publish(reply, handle.ID())
}

func (s *Server) channelPlay(ctx context.Context, reply string, req *proxy.Request) {
	if req.ChannelPlay.PlaybackID == "" {
		req.ChannelPlay.PlaybackID = uuid.NewV1().String()
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", req.ChannelPlay.ID)
		s.Dialog.Bind(req.Metadata.Dialog, "playback", req.ChannelPlay.PlaybackID)
	}

	ph, err := s.ari.Channel().Play(req.ChannelPlay.ID, req.ChannelPlay.PlaybackID, req.ChannelPlay.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: used to send nil
	s.nats.Publish(reply, ph.ID())
}

func (s *Server) channelRecord(ctx context.Context, reply string, req *proxy.Request) {
	if req.ChannelRecord.Name == "" {
		req.ChannelRecord.Name = uuid.NewV1().String()
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", req.ChannelRecord.ID)
		s.Dialog.Bind(req.Metadata.Dialog, "recording", req.ChannelRecord.Name)
	}

	lr, err := s.ari.Channel().Record(req.ChannelRecord.ID, req.ChannelRecord.Name, req.ChannelRecord.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, lr.ID())
}

func (s *Server) channelRing(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelRing.ID
	err := s.ari.Channel().Ring(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelSendDTMF(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelSendDTMF.ID
	dtmf := req.ChannelSendDTMF.DTMF
	opts := req.ChannelSendDTMF.Options
	err := s.ari.Channel().SendDTMF(id, dtmf, opts)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelSilence(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelSilence.ID
	err := s.ari.Channel().Silence(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelSnoop(ctx context.Context, reply string, req *proxy.Request) {
	ch, err := s.ari.Channel().Snoop(req.ChannelSnoop.ID, req.ChannelSnoop.SnoopID, req.ChannelSnoop.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", ch.ID())
	}

	s.nats.Publish(reply, ch.ID())
}

func (s *Server) channelStopHold(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel().StopHold(req.ChannelStopHold.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelStopMOH(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel().StopMOH(req.ChannelStopMOH.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelStopRing(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel().StopRing(req.ChannelStopRing.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelStopSilence(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelStopSilence.ID
	err := s.ari.Channel().StopSilence(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// check for existence
	_, err := s.ari.Channel().Data(req.ChannelSubscribe.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", req.ChannelSubscribe.ID)
	}

	s.sendError(reply, nil)
}

func (s *Server) channelUnmute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Unmute(req.ChannelUnmute.ID, req.ChannelUnmute.Direction))
}

func (s *Server) channelVariableGet(ctx context.Context, reply string, req *proxy.Request) {
	val, err := s.ari.Channel().Variables(req.ChannelVariables.ID).Get(req.ChannelVariables.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, val)
}

func (s *Server) channelVariableSet(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel().Variables(req.ChannelVariables.ID).Set(req.ChannelVariables.Name, req.ChannelVariables.Set.Value)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
