package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingStoredCopy(ctx context.Context, reply string, req *proxy.Request) {
	h, err := s.ari.StoredRecording().Copy(req.Key, req.RecordingStoredCopy.Destination)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) recordingStoredData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.StoredRecording().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			StoredRecording: data,
		},
	})
}

func (s *Server) recordingStoredGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.StoredRecording().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) recordingStoredDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.StoredRecording().Delete(req.Key))
}

func (s *Server) recordingStoredList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.StoredRecording().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Keys: list,
	})
}
