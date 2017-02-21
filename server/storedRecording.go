package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingStoredCopy(ctx context.Context, reply string, req *proxy.Request) {
	id := req.RecordingStoredCopy.ID
	dest := req.RecordingStoredCopy.Destination

	srd, err := s.ari.Recording.Stored.Copy(id, dest)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &srd)
}

func (s *Server) recordingStoredData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Recording.Stored.Data(req.RecordingStoredData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &data)
}

func (s *Server) recordingStoredDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Recording.Stored.Delete(req.RecordingStoredDelete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) recordingStoredList(ctx context.Context, reply string, req *proxy.Request) {
	handles, err := s.ari.Recording.Stored.List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var ret []string
	for _, h := range handles {
		ret = append(ret, h.ID())
	}

	s.nats.Publish(reply, &handles)
}
