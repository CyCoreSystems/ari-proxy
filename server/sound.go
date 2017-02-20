package ariproxy

/*
func (ins *Instance) sound() {
	ins.subscribe("ari.sounds.all", func(msg *session.Message, reply Reply) {

		var filters map[string]string
		if err := json.Unmarshal(msg.Payload, &filters); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		if len(filters) == 0 {
			filters = nil // just send nil to upstream if empty. makes tests easier
		}

		sx, err := ins.upstream.Sound.List(filters)
		if err != nil {
			reply(nil, err)
			return
		}

		var sounds []string
		for _, sound := range sx {
			sounds = append(sounds, sound.ID())
		}

		reply(sounds, nil)
	})

	ins.subscribe("ari.sounds.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		sd, err := ins.upstream.Sound.Data(name)
		reply(&sd, err)
	})

}
*/
