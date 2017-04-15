package server

import (
	"context"
	"errors"
	"strings"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) applicationData(ctx context.Context, reply string, req *proxy.Request) {
	app, err := s.ari.Application().Data(req.ApplicationData.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Metadata:    s.Metadata(req.Metadata.Dialog),
			Application: app,
		},
	})
}

func (s *Server) applicationList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Application().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	resp := proxy.EntityList{List: []*proxy.Entity{}}
	for _, i := range list {
		resp.List = append(resp.List, &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       i.ID(),
		})
	}

	s.nats.Publish(reply, &proxy.Response{
		EntityList: &resp,
	})
}

func (s *Server) applicationGet(ctx context.Context, reply string, req *proxy.Request) {
	app, err := s.ari.Application().Data(req.ApplicationGet.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Entity{
		Metadata: s.Metadata(req.Metadata.Dialog),
		ID:       app.Name,
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
	err := s.ari.Application().Subscribe(req.ApplicationSubscribe.Name, req.ApplicationSubscribe.EventSource)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	if req.Metadata.Dialog != "" {
		eType, eID, err := parseEventSource(req.ApplicationSubscribe.EventSource)
		if err != nil {
			s.Log.Warn("failed to parse event source", "error", err, "eventsource", req.ApplicationSubscribe.EventSource)
		} else {
			s.Dialog.Bind(req.Metadata.Dialog, eType, eID)
		}
	}

	s.sendError(reply, nil)
}

func (s *Server) applicationUnsubscribe(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Application().Unsubscribe(req.ApplicationUnsubscribe.Name, req.ApplicationUnsubscribe.EventSource)
	s.sendError(reply, err)
}
