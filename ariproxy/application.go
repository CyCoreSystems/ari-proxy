package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/session"
)

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
