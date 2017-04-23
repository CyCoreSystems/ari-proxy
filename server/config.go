package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskConfigData(ctx context.Context, reply string, req *proxy.Request) {

	data, err := s.ari.Asterisk().Config().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Config: data,
		},
	})
}

func (s *Server) asteriskConfigDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Config().Delete(req.Key))
}

func (s *Server) asteriskConfigUpdate(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Config().Update(req.Key, req.AsteriskConfig.Update.Tuples))
}
