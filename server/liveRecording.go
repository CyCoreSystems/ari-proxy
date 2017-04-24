package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingLiveData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.LiveRecording().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, proxy.Response{
		Data: &proxy.EntityData{
			LiveRecording: data,
		},
	})
}

func (s *Server) recordingLiveGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.LiveRecording().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) recordingLiveMute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Mute(req.Key))
}

func (s *Server) recordingLivePause(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Pause(req.Key))
}

func (s *Server) recordingLiveResume(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Resume(req.Key))
}

func (s *Server) recordingLiveScrap(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Scrap(req.Key))
}

func (s *Server) recordingLiveStop(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Stop(req.Key))
}

func (s *Server) recordingLiveUnmute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Unmute(req.Key))
}
