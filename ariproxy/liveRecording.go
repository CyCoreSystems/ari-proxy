package ariproxy

import "github.com/CyCoreSystems/ari-proxy/session"

func (ins *Instance) liveRecording() {
	ins.subscribe("ari.recording.live.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		lrd, err := ins.upstream.Recording.Live.Data(name)
		reply(lrd, err)
	})

	ins.subscribe("ari.recording.live.stop", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Stop(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.pause", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Pause(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.resume", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Resume(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.mute", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Mute(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.unmute", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Unmute(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Delete(name)
		reply(nil, err)
	})

	ins.subscribe("ari.recording.live.scrap", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Recording.Live.Scrap(name)
		reply(nil, err)
	})

}
