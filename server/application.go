package server

import (
	"context"
	"errors"
	"strings"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) applicationData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Application().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Application: data,
		},
	})
}

func (s *Server) applicationList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Application().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) applicationGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Application().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func parseEventSource(src string) (string, string, error) {
	var err error

	pieces := strings.Split(src, ":")
	if len(pieces) != 2 {
		return "", "", errors.New("Invalid EventSource")
	}

	switch pieces[0] {
	case "channel":
	case "bridge":
	case "endpoint":
	case "deviceState":
	default:
		err = errors.New("Unhandled EventSource type")
	}
	return pieces[0], pieces[1], err
}

func (s *Server) applicationSubscribe(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Application().Subscribe(req.Key, req.ApplicationSubscribe.EventSource)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Key.Dialog != "" {
		eType, eID, err := parseEventSource(req.ApplicationSubscribe.EventSource)
		if err != nil {
			s.Log.Warn("failed to parse event source", "error", err, "eventsource", req.ApplicationSubscribe.EventSource)
		} else {
			s.Dialog.Bind(req.Key.Dialog, eType, eID)
		}
	}

	s.sendError(reply, nil)
}

func (s *Server) applicationUnsubscribe(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Application().Unsubscribe(req.Key, req.ApplicationUnsubscribe.EventSource)
	s.sendError(reply, err)
}
