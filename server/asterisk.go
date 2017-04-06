package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskInfo(ctx context.Context, reply string, req *proxy.Request) {
	info, err := s.ari.Asterisk().Info("")
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.DataResponse{
		AsteriskInfo: info,
	})
}

func (s *Server) asteriskReloadModule(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().ReloadModule(req.AsteriskReloadModule.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) asteriskVariableGet(ctx context.Context, reply string, req *proxy.Request) {
	val, err := s.ari.Asterisk().Variables().Get(req.AsteriskVariables.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.DataResponse{
		Variable: val,
	})
}

func (s *Server) asteriskVariableSet(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Variables().Set(req.AsteriskVariables.Name, req.AsteriskVariables.Set.Value)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
