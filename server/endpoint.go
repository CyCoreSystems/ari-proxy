package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) endpointData(ctx context.Context, reply string, req *proxy.Request) {
	ed, err := s.ari.Endpoint.Data(req.EndpointData.Tech, req.EndpointData.Resource)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &ed)
}

func (s *Server) endpointList(ctx context.Context, reply string, req *proxy.Request) {
	ex, err := s.ari.Endpoint.List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &ex)
}

func (s *Server) endpointListByTech(ctx context.Context, reply string, req *proxy.Request) {
	ex, err := s.ari.Endpoint.ListByTech(req.EndpointListByTech.Tech)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &ex)
}
