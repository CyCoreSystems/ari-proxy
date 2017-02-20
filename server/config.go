package ariproxy

/*
func (ins *Instance) config() {
	ins.subscribe("ari.asterisk.config.data", func(msg *session.Message, reply Reply) {
		name := msg.Object

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		cd, err := ins.upstream.Asterisk.Config().Data(items[0], items[1], items[2])
		reply(&cd.Fields, err)
	})

	ins.subscribe("ari.asterisk.config.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		err := ins.upstream.Asterisk.Config().Delete(items[0], items[1], items[2])
		reply(nil, err)
	})

	ins.subscribe("ari.asterisk.config.update", func(msg *session.Message, reply Reply) {
		name := msg.Object

		items := strings.Split(name, ".")
		if len(items) != 3 {
			reply(nil, errors.New("Malformed config ID in request"))
			return
		}

		var fl []ari.ConfigTuple
		if err := json.Unmarshal(msg.Payload, &fl); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Asterisk.Config().Update(items[0], items[1], items[2], fl)
		reply(nil, err)
	})

}
*/
