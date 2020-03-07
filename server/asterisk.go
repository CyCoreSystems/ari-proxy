package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
)

func (s *Server) asteriskInfo(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Info(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Asterisk: data,
		},
	})
}

func (s *Server) asteriskVariableGet(ctx context.Context, reply string, req *proxy.Request) {
	val, err := s.ari.Asterisk().Variables().Get(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Variable: val,
		},
	})
}

func (s *Server) asteriskVariableSet(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Variables().Set(req.Key, req.AsteriskVariableSet.Value)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
