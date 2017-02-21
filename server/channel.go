package server

import (
	"context"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) channelAnswer(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelAnswer.ID
	err := s.ari.Channel.Answer(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelBusy(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel.Busy(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelCongestion(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel.Congestion(ID)
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

	handle, err := s.ari.Channel.Create(create)
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
	d, err := s.ari.Channel.Data(req.ChannelData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &d)
}

func (s *Server) channelContinue(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelContinue.ID

	cont := req.ChannelContinue.ContinueRequest

	err := s.ari.Channel.Continue(id, cont.Context, cont.Extension, cont.Priority)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelDial(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelDial.ID
	dial := req.ChannelDial.DialRequest

	//TODO: confirm time is in Seconds, the ARI documentation does not list it for Dial
	err := s.ari.Channel.Dial(id, dial.Caller, time.Duration(dial.Timeout)*time.Second)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelHangup(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelHangup.ID
	reason := req.ChannelHangup.Reason
	err := s.ari.Channel.Hangup(id, reason)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelHold(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelHold.ID
	err := s.ari.Channel.Hold(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelList(ctx context.Context, reply string, req *proxy.Request) {
	cx, err := s.ari.Channel.List()
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
	err := s.ari.Channel.MOH(id, music)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelMute(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelMute.ID
	dir := req.ChannelMute.Direction
	err := s.ari.Channel.Mute(id, dir)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelOriginate(ctx context.Context, reply string, req *proxy.Request) {
	orig := req.ChannelOriginate.OriginateRequest

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.ChannelID)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.OtherChannelID)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", orig.Originator)
	}

	handle, err := s.ari.Channel.Originate(orig)
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
	id := req.ChannelPlay.ID
	pr := req.ChannelPlay.PlayRequest
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", id)
		s.Dialog.Bind(req.Metadata.Dialog, "playback", pr.PlaybackID)
	}

	ph, err := s.ari.Channel.Play(id, pr.PlaybackID, pr.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "playback", ph.ID())
	}

	//NOTE: used to send nil
	s.nats.Publish(reply, ph.ID())
}

func (s *Server) channelRecord(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelRecord.ID
	rr := req.ChannelRecord.RecordRequest

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", id)
		s.Dialog.Bind(req.Metadata.Dialog, "recording", rr.Name)
	}

	var opts ari.RecordingOptions

	opts.Format = rr.Format
	opts.MaxDuration = time.Duration(rr.MaxDuration) * time.Second
	opts.MaxSilence = time.Duration(rr.MaxSilence) * time.Second
	opts.Exists = rr.IfExists
	opts.Beep = rr.Beep
	opts.Terminate = rr.TerminateOn

	lr, err := s.ari.Channel.Record(id, rr.Name, &opts)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "recording", lr.ID())
	}

	//NOTE: used to send nil
	s.nats.Publish(reply, lr.ID())
}

func (s *Server) channelRing(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelRing.ID
	err := s.ari.Channel.Ring(id)
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
	err := s.ari.Channel.SendDTMF(id, dtmf, opts)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelSilence(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelSilence.ID
	err := s.ari.Channel.Silence(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelSnoop(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelSnoop.ID
	snoop := req.ChannelSnoop.SnoopRequest

	if req.Metadata.Dialog != "" {
		//TODO: confirm that snoopID is a channel
		s.Dialog.Bind(req.Metadata.Dialog, "channel", snoop.SnoopID)
	}

	ch, err := s.ari.Channel.Snoop(id, snoop.SnoopID, snoop.App, snoop.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "channel", ch.ID())
	}

	//NOTE: this used to send nil
	s.nats.Publish(reply, ch.ID())
}

func (s *Server) channelStopHold(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel.StopHold(req.ChannelStopHold.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelStopMOH(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel.StopMOH(req.ChannelStopMOH.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelStopRing(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel.StopRing(req.ChannelStopRing.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelStopSilence(ctx context.Context, reply string, req *proxy.Request) {
	id := req.ChannelStopSilence.ID
	err := s.ari.Channel.StopSilence(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) channelSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// check for existence
	_, err := s.ari.Channel.Data(req.ChannelSubscribe.ID)
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

	err := s.ari.Channel.Unmute(req.ChannelUnmute.ID, req.ChannelUnmute.Direction)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelVariableGet(ctx context.Context, reply string, req *proxy.Request) {
	val, err := s.ari.Channel.Variables(req.ChannelVariables.ID).Get(req.ChannelVariables.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, val)
}

func (s *Server) channelVariableSet(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Channel.Variables(req.ChannelVariables.ID).Set(req.ChannelVariables.Name, req.ChannelVariables.Set.Value)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
