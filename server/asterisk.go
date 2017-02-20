package ariproxy

/*
func (ins *Instance) asterisk() {

	ins.subscribe("ari.asterisk.reload", func(msg *session.Message, reply Reply) {
		err := ins.upstream.Asterisk.ReloadModule(msg.Object)
		reply(nil, err)
	})

	ins.subscribe("ari.asterisk.info", func(_ *session.Message, reply Reply) {
		ai, err := ins.upstream.Asterisk.Info("")
		reply(ai, err)
	})

	ins.subscribe("ari.asterisk.variables.get", func(msg *session.Message, reply Reply) {
		val, err := ins.upstream.Asterisk.Variables().Get(msg.Object)
		reply(val, err)
	})

	ins.subscribe("ari.asterisk.variables.set", func(msg *session.Message, reply Reply) {
		var value string
		if err := json.Unmarshal(msg.Payload, &value); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Asterisk.Variables().Set(msg.Object, value)
		reply(nil, err)
	})

}
*/
