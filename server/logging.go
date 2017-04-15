package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskLoggingList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Asterisk().Logging().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var ret proxy.EntityList
	for _, i := range list {
		ret.List = append(ret.List, &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			Kind:     "logging",
			ID:       i.ID(),
		})
	}
	s.nats.Publish(reply, &proxy.Response{EntityList: &ret})
}

func (s *Server) asteriskLoggingCreate(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Logging().Create(req.AsteriskLogging.Create.ID, req.AsteriskLogging.Create.Config)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskLoggingData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Logging().Data(req.AsteriskLogging.Data.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Log: data,
		},
	})
}

func (s *Server) asteriskLoggingRotate(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Logging().Rotate(req.AsteriskLogging.Rotate.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskLoggingDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Logging().Delete(req.AsteriskLogging.Delete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}
