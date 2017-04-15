package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) recordingStoredCopy(ctx context.Context, reply string, req *proxy.Request) {
	id := req.RecordingStoredCopy.ID
	dest := req.RecordingStoredCopy.Destination

	srd, err := s.ari.StoredRecording().Copy(id, dest)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Entity: &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       srd.ID(),
		},
	})
}

func (s *Server) recordingStoredData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.StoredRecording().Data(req.RecordingStoredData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Metadata:        s.Metadata(req.Metadata.Dialog),
			StoredRecording: data,
		},
	})
}

func (s *Server) recordingStoredDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.StoredRecording().Delete(req.RecordingStoredDelete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) recordingStoredList(ctx context.Context, reply string, req *proxy.Request) {
	handles, err := s.ari.StoredRecording().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var el proxy.EntityList
	for _, sr := range handles {
		el.List = append(el.List, &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       sr.ID(),
		})
	}

	s.nats.Publish(reply, &proxy.Response{
		EntityList: &el,
	})
}
