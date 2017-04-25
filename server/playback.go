package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) playbackControl(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Playback().Control(req.Key, req.PlaybackControl.Command))
}

func (s *Server) playbackData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Playback().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Playback: data,
		},
	})
}

func (s *Server) playbackGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Playback().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) playbackStop(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Playback().Stop(req.Key))
}

func (s *Server) playbackSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// bind dialog
	if req.Key.Dialog != "" {
		s.Dialog.Bind(req.Key.Dialog, "playback", req.Key.ID)
	}

	s.sendError(reply, nil)
}
