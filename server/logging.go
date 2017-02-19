package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/session"
)

func (ins *Instance) logging() {
	ins.subscribe("ari.logging.create", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var config string
		if err := json.Unmarshal(msg.Payload, &config); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Asterisk.Logging().Create(name, config)
		reply(nil, err)
		return
	})

	ins.subscribe("ari.logging.all", func(msg *session.Message, reply Reply) {
		ld, err := ins.upstream.Asterisk.Logging().List()
		reply(ld, err)
	})

	ins.subscribe("ari.logging.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Asterisk.Logging().Delete(name)
		reply(nil, err)
	})

	ins.subscribe("ari.logging.rotate", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Asterisk.Logging().Rotate(name)
		reply(nil, err)
	})

}
