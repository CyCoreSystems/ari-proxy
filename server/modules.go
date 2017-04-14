package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskModuleLoad(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Modules().Load(req.AsteriskModules.Load.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskModuleUnload(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Modules().Unload(req.AsteriskModules.Unload.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskModuleReload(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Asterisk().Modules().Reload(req.AsteriskModules.Reload.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.sendError(reply, nil)
}

func (s *Server) asteriskModuleData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Modules().Data(req.AsteriskModules.Data.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.DataResponse{
		ModuleData: data,
	})
}

func (s *Server) asteriskModuleList(ctx context.Context, reply string, req *proxy.Request) {
	mx, err := s.ari.Asterisk().Modules().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var el proxy.EntityList
	for _, m := range mx {
		el.List = append(el.List, &proxy.Entity{
			ID: m.ID(),
		})
	}

	s.nats.Publish(reply, &proxy.Response{
		EntityList: &el,
	})
}
