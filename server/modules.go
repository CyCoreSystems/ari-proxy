package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskModuleLoad(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Modules().Load(req.Key))
}

func (s *Server) asteriskModuleUnload(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Modules().Unload(req.Key))
}

func (s *Server) asteriskModuleReload(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Modules().Reload(req.Key))
}

func (s *Server) asteriskModuleData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Modules().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Module: data,
		},
	})
}

func (s *Server) asteriskModuleGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Modules().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) asteriskModuleList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Asterisk().Modules().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Keys: list,
	})
}
