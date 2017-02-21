package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingLiveData(ctx context.Context, reply string, req *proxy.Request) {
	lrd, err := s.ari.Recording.Live.Data(req.RecordingLiveData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &lrd)
}

func (s *Server) recordingLiveDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Delete(req.RecordingLiveDelete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLiveMute(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Mute(req.RecordingLiveMute.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLivePause(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Pause(req.RecordingLivePause.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLiveResume(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Resume(req.RecordingLiveResume.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLiveScrap(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Scrap(req.RecordingLiveScrap.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLiveStop(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Stop(req.RecordingLiveStop.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) recordingLiveUnmute(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Live.Unmute(req.RecordingLiveUnmute.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)

}
