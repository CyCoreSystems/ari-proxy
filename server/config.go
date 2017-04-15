package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) asteriskConfigData(ctx context.Context, reply string, req *proxy.Request) {

	cd, err := s.ari.Asterisk().Config().Data(
		req.AsteriskConfig.ConfigClass,
		req.AsteriskConfig.ObjectType,
		req.AsteriskConfig.ID,
	)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Metadata: s.Metadata(req.Metadata.Dialog),
			Config:   cd,
		},
	})
}

func (s *Server) asteriskConfigDelete(ctx context.Context, reply string, req *proxy.Request) {

	err := s.ari.Asterisk().Config().Delete(
		req.AsteriskConfig.ConfigClass,
		req.AsteriskConfig.ObjectType,
		req.AsteriskConfig.ID,
	)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) asteriskConfigUpdate(ctx context.Context, reply string, req *proxy.Request) {

	err := s.ari.Asterisk().Config().Update(
		req.AsteriskConfig.ConfigClass,
		req.AsteriskConfig.ObjectType,
		req.AsteriskConfig.ID,
		req.AsteriskConfig.Update.Tuples,
	)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
