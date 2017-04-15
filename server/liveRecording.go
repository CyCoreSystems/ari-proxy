package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingLiveData(ctx context.Context, reply string, req *proxy.Request) {
	lrd, err := s.ari.LiveRecording().Data(req.RecordingLiveData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, proxy.Response{
		Data: &proxy.EntityData{
			LiveRecording: lrd,
		},
	})
}

func (s *Server) recordingLiveDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Delete(req.RecordingLiveDelete.ID))
}

func (s *Server) recordingLiveMute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Mute(req.RecordingLiveMute.ID))
}

func (s *Server) recordingLivePause(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Pause(req.RecordingLivePause.ID))
}

func (s *Server) recordingLiveResume(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Resume(req.RecordingLiveResume.ID))
}

func (s *Server) recordingLiveScrap(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Scrap(req.RecordingLiveScrap.ID))
}

func (s *Server) recordingLiveStop(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Stop(req.RecordingLiveStop.ID))
}

func (s *Server) recordingLiveUnmute(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.LiveRecording().Unmute(req.RecordingLiveUnmute.ID))
}
