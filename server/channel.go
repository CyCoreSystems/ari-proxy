package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) channelAnswer(ctx context.Context, reply string, req *proxy.Request) {
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
	}

	s.sendError(reply, s.ari.Channel().Answer(req.Key))
}

func (s *Server) channelBusy(ctx context.Context, reply string, req *proxy.Request) {
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
	}

	s.sendError(reply, s.ari.Channel().Busy(req.Key))
}

func (s *Server) channelCongestion(ctx context.Context, reply string, req *proxy.Request) {
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
	}

	s.sendError(reply, s.ari.Channel().Congestion(req.Key))
}

func (s *Server) channelCreate(ctx context.Context, reply string, req *proxy.Request) {
	create := req.ChannelCreate.ChannelCreateRequest

	if create.ChannelID == "" {
		create.ChannelID = uuid.NewV1().String()
	}

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", create.ChannelID)
		s.Dialog.Bind(req.Key.Dialog, "channel", create.OtherChannelID)
	}

	h, err := s.ari.Channel().Create(create)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) channelData(ctx context.Context, reply string, req *proxy.Request) {
	d, err := s.ari.Channel().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Channel: d,
		},
	})
}

func (s *Server) channelContinue(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Continue(req.Key, req.ChannelContinue.Context, req.ChannelContinue.Extension, req.ChannelContinue.Priority))
}

func (s *Server) channelDial(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Dial(req.Key, req.ChannelDial.Caller, req.ChannelDial.Timeout))
}

func (s *Server) channelHangup(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Hangup(req.Key, req.ChannelHangup.Reason))
}

func (s *Server) channelHold(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Hold(req.Key))
}

func (s *Server) channelList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Channel().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) channelMOH(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().MOH(req.Key, req.ChannelMOH.Music))
}

func (s *Server) channelMute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Mute(req.Key, req.ChannelMute.Direction))
}

func (s *Server) channelOriginate(ctx context.Context, reply string, req *proxy.Request) {
	orig := req.ChannelOriginate.OriginateRequest

	if orig.ChannelID == "" {
		orig.ChannelID = uuid.NewV1().String()
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", orig.ChannelID)
		s.Dialog.Bind(req.Key.Dialog, "channel", orig.OtherChannelID)
		s.Dialog.Bind(req.Key.Dialog, "channel", orig.Originator)
	}

	h, err := s.ari.Channel().Originate(orig)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) channelPlay(ctx context.Context, reply string, req *proxy.Request) {
	if req.ChannelPlay.PlaybackID == "" {
		req.ChannelPlay.PlaybackID = uuid.NewV1().String()
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "playback", req.ChannelPlay.PlaybackID)
	}

	ph, err := s.ari.Channel().Play(req.Key, req.ChannelPlay.PlaybackID, req.ChannelPlay.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: ph.Key(),
	})
}

func (s *Server) channelRecord(ctx context.Context, reply string, req *proxy.Request) {
	if req.ChannelRecord.Name == "" {
		req.ChannelRecord.Name = uuid.NewV1().String()
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "recording", req.ChannelRecord.Name)
	}

	h, err := s.ari.Channel().Record(req.Key, req.ChannelRecord.Name, req.ChannelRecord.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) channelRing(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Ring(req.Key))
}

func (s *Server) channelSendDTMF(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().SendDTMF(req.Key, req.ChannelSendDTMF.DTMF, req.ChannelSendDTMF.Options))
}

func (s *Server) channelSilence(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Silence(req.Key))
}

func (s *Server) channelSnoop(ctx context.Context, reply string, req *proxy.Request) {
	if req.ChannelSnoop.SnoopID == "" {
		req.ChannelSnoop.SnoopID = uuid.NewV1().String()
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.ChannelSnoop.SnoopID)
	}

	h, err := s.ari.Channel().Snoop(req.Key, req.ChannelSnoop.SnoopID, req.ChannelSnoop.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) channelStopHold(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().StopHold(req.Key))
}

func (s *Server) channelStopMOH(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().StopMOH(req.Key))
}

func (s *Server) channelStopRing(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().StopRing(req.Key))
}

func (s *Server) channelStopSilence(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().StopSilence(req.Key))
}

func (s *Server) channelSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "channel", req.Key.ID)
	}

	s.sendError(reply, nil)
}

func (s *Server) channelUnmute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Unmute(req.Key, req.ChannelUnmute.Direction))
}

func (s *Server) channelVariableGet(ctx context.Context, reply string, req *proxy.Request) {
	val, err := s.ari.Channel().Variables(req.Key).Get(req.ChannelVariables.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	// TODO: return needs metadata/key data somehow
	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Variable: val,
		},
	})
}

func (s *Server) channelVariableSet(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Channel().Variables(req.Key).Set(req.ChannelVariables.Name, req.ChannelVariables.Set.Value))
}
