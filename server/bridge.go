package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
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

	bh, err := s.ari.Bridge().Create(req.Key, req.BridgeCreate.Type, req.BridgeCreate.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: bh.Key(),
	})
}

func (s *Server) bridgeData(ctx context.Context, reply string, req *proxy.Request) {
	bd, err := s.ari.Bridge().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
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

func (s *Server) bridgeList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Bridge().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Keys: list,
	})
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

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &proxy.Response{
		Key: ph.Key(),
	})
}

func (s *Server) bridgeRecord(ctx context.Context, reply string, req *proxy.Request) {

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

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
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
