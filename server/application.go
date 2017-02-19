package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari-proxy/session"
)

func (s *Server) applicationData(reply string, req *proxy.Request) {
	app, err := s.ari.Application.Data(req.ApplicationData.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	if app == nil {
		s.sendNotFound(reply)
		return
	}

	s.nats.Publish(reply, &app)
}

func (s *Server) applicationList(reply string, req *proxy.Request) {
	list, err := s.ari.Application.List()
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
	s.nats.Publish(reply, &resp)
}

func (s *Server) applicationGet(reply string, req *proxy.Request) {
	app, err := s.ari.Application.Get(req.ApplicationGet.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	if app == nil {
		s.sendNotFound(reply)
		return
	}

	s.nats.Publish(reply, &proxy.Entity{
		Metadata: s.Metadata(req.Metadata.Dialog),
		ID:       app.ID(),
	})
}

func (ins *Instance) application() {
	ins.subscribe("ari.applications.all", func(_ *session.Message, reply Reply) {
		ax, err := ins.upstream.Application.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var apps []string
		for _, a := range ax {
			apps = append(apps, a.ID())
		}

		reply(apps, nil)
	})

	ins.subscribe("ari.applications.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		data, err := ins.upstream.Application.Data(name)
		reply(data, err)
	})

	ins.subscribe("ari.applications.subscribe", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var eventSource string
		if err := json.Unmarshal(msg.Payload, &eventSource); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Application.Subscribe(name, eventSource)
		reply(nil, err)
	})

	ins.subscribe("ari.applications.unsubscribe", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var eventSource string
		if err := json.Unmarshal(msg.Payload, &eventSource); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Application.Unsubscribe(name, eventSource)
		reply(nil, err)
	})
}
