package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/session"
)

func (ins *Instance) device() {
	ins.subscribe("ari.devices.all", func(msg *session.Message, reply Reply) {

		dx, err := ins.upstream.DeviceState.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var ret []string
		for _, device := range dx {
			ret = append(ret, device.ID())
		}

		reply(ret, nil)
	})

	ins.subscribe("ari.devices.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		dd, err := ins.upstream.DeviceState.Data(name)
		reply(dd, err)
	})

	ins.subscribe("ari.devices.update", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var state string
		if err := json.Unmarshal(msg.Payload, &state); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.DeviceState.Update(name, state)
		reply(nil, err)
	})

	ins.subscribe("ari.devices.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.DeviceState.Delete(name)
		reply(nil, err)
	})

}
