package ariproxy

/*
func (ins *Instance) bridge() {

	ins.subscribe("ari.bridges.all", func(msg *session.Message, reply Reply) {

		bx, err := ins.upstream.Bridge.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var bridges []string
		for _, bridge := range bx {
			bridges = append(bridges, bridge.ID())
		}

		reply(bridges, nil)
	})

	ins.subscribe("ari.bridges.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		bd, err := ins.upstream.Bridge.Data(name)
		reply(&bd, err)
		return
	})

	ins.subscribe("ari.bridges.create", func(msg *session.Message, reply Reply) {

		var req client.CreateBridgeRequest
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		bh, err := ins.upstream.Bridge.Create(req.ID, req.Type, req.Name)

		if err != nil {
			reply(nil, err)
			return
		}

		ins.server.cache.Add(req.ID, ins)

		reply(bh.ID(), err)
	})

	ins.subscribe("ari.bridges.addChannel", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var channelID string
		if err := json.Unmarshal(msg.Payload, &channelID); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Bridge.AddChannel(name, channelID)
		reply(nil, err)
	})

	ins.subscribe("ari.bridges.removeChannel", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var channelID string
		if err := json.Unmarshal(msg.Payload, &channelID); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Bridge.RemoveChannel(name, channelID)
		reply(nil, err)
	})

	ins.subscribe("ari.bridges.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Bridge.Delete(name)
		reply(nil, err)
	})

	ins.subscribe("ari.bridges.play", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var pr client.PlayRequest
		if err := json.Unmarshal(msg.Payload, &pr); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		_, err := ins.upstream.Bridge.Play(name, pr.PlaybackID, pr.MediaURI)
		reply(nil, err)
	})

	ins.subscribe("ari.bridges.record", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var rr client.RecordRequest
		if err := json.Unmarshal(msg.Payload, &rr); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		ins.server.cache.Add(rr.Name, ins)

		var opts ari.RecordingOptions

		opts.Format = rr.Format
		opts.MaxDuration = time.Duration(rr.MaxDuration) * time.Second
		opts.MaxSilence = time.Duration(rr.MaxSilence) * time.Second
		opts.Exists = rr.IfExists
		opts.Beep = rr.Beep
		opts.Terminate = rr.TerminateOn

		_, err := ins.upstream.Bridge.Record(name, rr.Name, &opts)
		reply(nil, err)
	})
}
*/
