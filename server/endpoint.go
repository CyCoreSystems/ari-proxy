package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
)

func (s *Server) endpointData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Endpoint().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Endpoint: data,
		},
	})
}

func (s *Server) endpointGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Endpoint().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) endpointList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Endpoint().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) endpointListByTech(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Endpoint().ListByTech(req.EndpointListByTech.Tech, req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}
