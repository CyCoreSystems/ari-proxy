package ariproxy

/*
func (ins *Instance) modules() {
	ins.subscribe("ari.modules.all", func(msg *session.Message, reply Reply) {
		mx, err := ins.upstream.Asterisk.Modules().List()
		if err != nil {
			reply(nil, err)
			return
		}

		var modules []string
		for _, m := range mx {
			modules = append(modules, m.ID())
		}

		reply(modules, nil)
	})

	ins.subscribe("ari.modules.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		data, err := ins.upstream.Asterisk.Modules().Data(name)
		reply(data, err)
	})

	ins.subscribe("ari.modules.load", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Asterisk.Modules().Load(name)
		reply(nil, err)
	})

	ins.subscribe("ari.modules.unload", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Asterisk.Modules().Unload(name)
		reply(nil, err)
	})

	ins.subscribe("ari.modules.reload", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Asterisk.Modules().Reload(name)
		reply(nil, err)
	})

}
*/
