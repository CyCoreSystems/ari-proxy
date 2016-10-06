package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/session"
)

func (ins *Instance) playback() {

	ins.subscribe("ari.playback.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		d, err := ins.upstream.Playback.Data(name)
		reply(&d, err)
	})

	ins.subscribe("ari.playback.control", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var command string
		if err := json.Unmarshal(msg.Payload, &command); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Playback.Control(name, command)
		reply(nil, err)
	})

	ins.subscribe("ari.playback.stop", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Playback.Stop(name)
		reply(nil, err)
	})

}
