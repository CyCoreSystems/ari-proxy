package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) playbackControl(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Playback.Control(req.PlaybackControl.ID, req.PlaybackControl.Command)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) playbackData(ctx context.Context, reply string, req *proxy.Request) {
	d, err := s.ari.Playback.Data(req.PlaybackData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &d)
}

func (s *Server) playbackStop(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Playback.Stop(req.PlaybackStop.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) playbackSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// check for existence
	_, err := s.ari.Playback.Data(req.PlaybackSubscribe.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "playback", req.PlaybackSubscribe.ID)
	}

	s.sendError(reply, nil)
}
