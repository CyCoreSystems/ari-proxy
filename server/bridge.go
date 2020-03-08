package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/rid"
)

func (s *Server) bridgeAddChannel(ctx context.Context, reply string, req *proxy.Request) {
	channel := req.BridgeAddChannel.Channel

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "channel", channel)
	}

	err := s.ari.Bridge().AddChannel(req.Key, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeCreate(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	h, err := s.ari.Bridge().Create(req.Key, req.BridgeCreate.Type, req.BridgeCreate.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) bridgeStageCreate(ctx context.Context, reply string, req *proxy.Request) {
	bh := s.ari.Bridge().Get(req.Key)

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	s.publish(reply, &proxy.Response{
		Key: bh.Key(),
	})
}

func (s *Server) bridgeData(ctx context.Context, reply string, req *proxy.Request) {
	bd, err := s.ari.Bridge().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Bridge: bd,
		},
	})
}

func (s *Server) bridgeDelete(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	err := s.ari.Bridge().Delete(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Bridge().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) bridgeList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Bridge().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) bridgeMOH(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	s.sendError(
		reply,
		s.ari.Bridge().MOH(req.Key, req.BridgeMOH.Class),
	)
}

func (s *Server) bridgeStopMOH(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	s.sendError(
		reply,
		s.ari.Bridge().StopMOH(req.Key),
	)
}

func (s *Server) bridgePlay(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "playback", req.BridgePlay.PlaybackID)
	}

	ph, err := s.ari.Bridge().Play(req.Key, req.BridgePlay.PlaybackID, req.BridgePlay.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: ph.Key(),
	})
}

func (s *Server) bridgeStagePlay(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Bridge().Data(req.Key)
	if err != nil || data == nil {
		s.sendError(reply, err)
		return
	}

	if req.BridgePlay.PlaybackID == "" {
		req.BridgePlay.PlaybackID = rid.New(rid.Playback)
	}

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "playback", req.BridgePlay.PlaybackID)
	}

	s.publish(reply, &proxy.Response{
		Key: s.ari.Playback().Get(ari.NewKey(ari.PlaybackKey, req.BridgePlay.PlaybackID)).Key(),
	})
}

func (s *Server) bridgeRecord(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Bridge().Data(req.Key)
	if err != nil || data == nil {
		s.sendError(reply, err)
		return
	}

	if req.BridgeRecord.Name == "" {
		req.BridgeRecord.Name = rid.New(rid.Recording)
	}

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "recording", req.BridgeRecord.Name)
	}

	h, err := s.ari.Bridge().Record(req.Key, req.BridgeRecord.Name, req.BridgeRecord.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) bridgeStageRecord(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Bridge().Data(req.Key)
	if err != nil || data == nil {
		s.sendError(reply, err)
		return
	}

	if req.BridgeRecord.Name == "" {
		req.BridgeRecord.Name = rid.New(rid.Recording)
	}

	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", data.ID)
		s.Dialog.Bind(req.Key.Dialog, "recording", req.BridgeRecord.Name)
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key.New(ari.LiveRecordingKey, req.BridgeRecord.Name),
	})
}

func (s *Server) bridgeRemoveChannel(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "channel", req.BridgeRemoveChannel.Channel)
	}

	err := s.ari.Bridge().RemoveChannel(req.Key, req.BridgeRemoveChannel.Channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeSubscribe(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeUnsubscribe(ctx context.Context, reply string, req *proxy.Request) {
	// no-op for now; may want to eventually optimize away the dialog subscription
	s.sendError(reply, nil)
}

func (s *Server) bridgeVideoSource(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
		s.Dialog.Bind(req.Key.Dialog, "channel", req.BridgeVideoSource.Channel)
	}

	err := s.ari.Bridge().VideoSource(req.Key, req.BridgeVideoSource.Channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeVideoSourceDelete(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "bridge", req.Key.ID)
	}

	err := s.ari.Bridge().VideoSourceDelete(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
