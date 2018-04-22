package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) deviceStateData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.DeviceState().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			DeviceState: data,
		},
	})
}

func (s *Server) deviceStateGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.DeviceState().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) deviceStateDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.DeviceState().Delete(req.Key))
}

func (s *Server) deviceStateList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.DeviceState().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) deviceStateUpdate(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.DeviceState().Update(req.Key, req.DeviceStateUpdate.State))
}
