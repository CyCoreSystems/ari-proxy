package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
)

func (s *Server) asteriskLoggingList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Asterisk().Logging().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) asteriskLoggingGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Logging().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) asteriskLoggingCreate(ctx context.Context, reply string, req *proxy.Request) {
	h, err := s.ari.Asterisk().Logging().Create(req.Key, req.AsteriskLoggingChannel.Levels)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: h.Key(),
	})
}

func (s *Server) asteriskLoggingData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Asterisk().Logging().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Log: data,
		},
	})
}

func (s *Server) asteriskLoggingRotate(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Logging().Rotate(req.Key))
}

func (s *Server) asteriskLoggingDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Asterisk().Logging().Delete(req.Key))
}
