package ariproxy

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) channelAnswer(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelAnswer.ID
	err := s.ari.Channel.Answer(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelBusy(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel.Busy(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelCongestion(ctx context.Context, reply string, req *proxy.Request) {
	ID := req.ChannelBusy.ID
	err := s.ari.Channel.Congestion(ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) channelCreate(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelDataContinue(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelDial(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelHangup(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelHold(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelList(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelMOH(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelMute(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelOriginate(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelPlay(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelRecord(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channeRing(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelSendDTMF(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelSilence(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelSnoop(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelStopHold(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelStopMOH(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelStopRing(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelStopSilence(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelSubscribe(ctx context.Context, reply string, req *proxy.Request) {

}

func (s *Server) channelUnmute(ctx context.Context, reply string, req *proxy.Request) {

}

/*
func (ins *Instance) channel() {

	ins.subscribe("ari.channels.all", func(msg *session.Message, reply Reply) {
		cx, err := ins.upstream.Channel.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var channels []string
		for _, channel := range cx {
			channels = append(channels, channel.ID())
		}

		reply(channels, nil)
	})

	ins.subscribe("ari.channels.originate", func(msg *session.Message, reply Reply) {
		var req ari.OriginateRequest

		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		handle, err := ins.upstream.Channel.Originate(req)

		if err != nil {
			reply(nil, err)
			return
		}

		reply(handle.ID(), nil)
	})

	ins.subscribe("ari.channels.create", func(msg *session.Message, reply Reply) {
		var req ari.ChannelCreateRequest

		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		if req.ChannelID != "" {
			ins.server.cache.Add(req.ChannelID, ins)
		}
		if req.OtherChannelID != "" {
			ins.server.cache.Add(req.OtherChannelID, ins)
		}

		handle, err := ins.upstream.Channel.Create(req)

		if err != nil {
			reply(nil, err)
			return
		}

		reply(handle.ID(), nil)
	})

	ins.subscribe("ari.channels.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		d, err := ins.upstream.Channel.Data(name)
		reply(&d, err)
	})

	ins.subscribe("ari.channels.answer", func(msg *session.Message, reply Reply) {
		name := msg.Object
		ins.log.Debug("answering channel", "msg.Command", msg.Command)
		err := ins.upstream.Channel.Answer(name)
		ins.log.Debug("answered channel", "msg.Command", msg.Command, "name", name, "error", err)

		reply(nil, err)
	})

	ins.subscribe("ari.channels.hangup", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var reason string
		if err := json.Unmarshal(msg.Payload, &reason); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Channel.Hangup(name, reason)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.ring", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.Ring(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.stopring", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.StopRing(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.hold", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.Hold(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.stophold", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.StopHold(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.mute", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var dir string
		if err := json.Unmarshal(msg.Payload, &dir); err != nil {
			reply(nil, &decodingError{msg.Command, err})
		}

		err := ins.upstream.Channel.Mute(name, dir)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.unmute", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var dir string
		if err := json.Unmarshal(msg.Payload, &dir); err != nil {
			reply(nil, &decodingError{msg.Command, err})
		}

		err := ins.upstream.Channel.Unmute(name, dir)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.silence", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.Silence(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.stopsilence", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.StopSilence(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.moh", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var music string
		if err := json.Unmarshal(msg.Payload, &music); err != nil {
			reply(nil, &decodingError{msg.Command, err})
		}

		err := ins.upstream.Channel.MOH(name, music)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.stopmoh", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Channel.StopMOH(name)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.play", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var pr client.PlayRequest
		if err := json.Unmarshal(msg.Payload, &pr); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		_, err := ins.upstream.Channel.Play(name, pr.PlaybackID, pr.MediaURI)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.dtmf", func(msg *session.Message, reply Reply) {
		name := msg.Object

		type request struct {
			Dtmf string           `json:"dtmf,omitempty"`
			Opts *ari.DTMFOptions `json:"options,omitempty"`
		}

		var req request
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Channel.SendDTMF(name, req.Dtmf, req.Opts)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.continue", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var req client.ContinueRequest
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Channel.Continue(name, req.Context, req.Extension, req.Priority)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.dial", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var req client.DialRequest
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		//TODO: confirm time is in Seconds, the ARI documentation does not list it for Dial
		err := ins.upstream.Channel.Dial(name, req.Caller, time.Duration(req.Timeout)*time.Second)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.snoop", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var req client.SnoopRequest
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		_, err := ins.upstream.Channel.Snoop(name, req.SnoopID, req.App, req.Options)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.record", func(msg *session.Message, reply Reply) {
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

		_, err := ins.upstream.Channel.Record(name, rr.Name, &opts)
		reply(nil, err)
	})

	ins.subscribe("ari.channels.variables.get", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var req client.GetChannelVariable
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		val, err := ins.upstream.Channel.Variables(name).Get(req.Name)
		reply(val, err)
	})

	ins.subscribe("ari.channels.variables.set", func(msg *session.Message, reply Reply) {
		name := msg.Object

		var req client.SetChannelVariable
		if err := json.Unmarshal(msg.Payload, &req); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Channel.Variables(name).Set(req.Name, req.Value)
		reply(nil, err)
	})

}
*/
