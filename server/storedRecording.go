package ariproxy

/*
func (ins *Instance) storedRecording() {
	ins.subscribe("ari.recording.stored.all", func(msg *session.Message, reply Reply) {
		handles, err := ins.upstream.Recording.Stored.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var ret []string
		for _, h := range handles {
			ret = append(ret, h.ID())
		}

		reply(ret, nil)
	})

	ins.subscribe("ari.recording.stored.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		srd, err := ins.upstream.Recording.Stored.Data(name)
		reply(srd, err)
	})

	ins.subscribe("ari.recording.stored.copy", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var dest string
		if err := json.Unmarshal(msg.Payload, &dest); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		srd, err := ins.upstream.Recording.Stored.Copy(name, dest)
		reply(srd, err)
	})

	ins.subscribe("ari.recording.stored.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Stored.Delete(name)
		reply(nil, err)
	})

}
*/
