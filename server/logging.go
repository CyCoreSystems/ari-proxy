package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskLoggingList(ctx context.Context, reply string, req *proxy.Request) {
	ld, err := s.ari.Asterisk.Logging().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &ld)
}

func (s *Server) asteriskLoggingCreate(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk.Logging().Create(req.AsteriskLogging.Create.ID, req.AsteriskLogging.Create.Config)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskLoggingRotate(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk.Logging().Rotate(req.AsteriskLogging.Rotate.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskLoggingDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk.Logging().Delete(req.AsteriskLogging.Delete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}
