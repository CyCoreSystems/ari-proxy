package client

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	uuid "github.com/satori/go.uuid"
)

type channel struct {
	c *Client
}

func (c *channel) Get(key *ari.Key) *ari.ChannelHandle {
	k, err := c.c.getRequest(&proxy.Request{
		Kind: "ChannelGet",
		Key:  key,
	})
	if err != nil {
		c.c.log.Warn("failed to make data request for channel", "error", err)
		return ari.NewChannelHandle(key, c, nil)
	}
	return ari.NewChannelHandle(k, c, nil)
}

func (c *channel) List(filter *ari.Key) ([]*ari.Key, error) {
	return c.c.listRequest(&proxy.Request{
		Kind: "ChannelList",
		Key:  filter,
	})
}

func (c *channel) Originate(key *ari.Key, o ari.OriginateRequest) (*ari.ChannelHandle, error) {
	k, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelOriginate",
		Key:  key,
		ChannelOriginate: &proxy.ChannelOriginate{
			OriginateRequest: o,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewChannelHandle(k.New(ari.ChannelKey, o.ChannelID), c, nil), nil
}

func (c *channel) StageOriginate(key *ari.Key, o ari.OriginateRequest) (*ari.ChannelHandle, error) {
	k, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelStageOriginate",
		Key:  key,
		ChannelOriginate: &proxy.ChannelOriginate{
			OriginateRequest: o,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewChannelHandle(k.New(ari.ChannelKey, o.ChannelID), c, func(h *ari.ChannelHandle) error {
		_, err := c.Originate(k.New(ari.ChannelKey, key.ID), o)
		return err
	}), nil
}

func (c *channel) Create(key *ari.Key, o ari.ChannelCreateRequest) (*ari.ChannelHandle, error) {
	k, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelCreate",
		Key:  key,
		ChannelCreate: &proxy.ChannelCreate{
			ChannelCreateRequest: o,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewChannelHandle(k.New(ari.ChannelKey, o.ChannelID), c, nil), nil
}

func (c *channel) Data(key *ari.Key) (*ari.ChannelData, error) {
	data, err := c.c.dataRequest(&proxy.Request{
		Kind: "ChannelData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Channel, nil
}

func (c *channel) Continue(key *ari.Key, context string, extension string, priority int) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelContinue",
		Key:  key,
		ChannelContinue: &proxy.ChannelContinue{
			Context:   context,
			Extension: extension,
			Priority:  priority,
		},
	})
}

func (c *channel) Busy(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelBusy",
		Key:  key,
	})
}

func (c *channel) Congestion(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelCongestion",
		Key:  key,
	})
}

func (c *channel) Hangup(key *ari.Key, reason string) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelHangup",
		Key:  key,
		ChannelHangup: &proxy.ChannelHangup{
			Reason: reason,
		},
	})
}

func (c *channel) Answer(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelAnswer",
		Key:  key,
	})
}

func (c *channel) Ring(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelRing",
		Key:  key,
	})
}

func (c *channel) StopRing(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelStopRing",
		Key:  key,
	})
}

func (c *channel) SendDTMF(key *ari.Key, dtmf string, opts *ari.DTMFOptions) error {
	if opts == nil {
		opts = &ari.DTMFOptions{}
	}
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelSendDTMF",
		Key:  key,
		ChannelSendDTMF: &proxy.ChannelSendDTMF{
			DTMF:    dtmf,
			Options: opts,
		},
	})
}

func (c *channel) Hold(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelHold",
		Key:  key,
	})
}

func (c *channel) StopHold(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelStopHold",
		Key:  key,
	})
}

func (c *channel) Mute(key *ari.Key, dir ari.Direction) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelMute",
		Key:  key,
		ChannelMute: &proxy.ChannelMute{
			Direction: dir,
		},
	})
}

func (c *channel) Unmute(key *ari.Key, dir ari.Direction) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelUnmute",
		Key:  key,
		ChannelMute: &proxy.ChannelMute{
			Direction: dir,
		},
	})
}

func (c *channel) MOH(key *ari.Key, moh string) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelMOH",
		Key:  key,
		ChannelMOH: &proxy.ChannelMOH{
			Music: moh,
		},
	})
}

func (c *channel) StopMOH(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelStopMOH",
		Key:  key,
	})
}

func (c *channel) Silence(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelSilence",
		Key:  key,
	})
}

func (c *channel) StopSilence(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelStopSilence",
		Key:  key,
	})
}

func (c *channel) Snoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) (*ari.ChannelHandle, error) {
	k, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelSnoop",
		Key:  key,
		ChannelSnoop: &proxy.ChannelSnoop{
			SnoopID: snoopID,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewChannelHandle(k.New(ari.ChannelKey, snoopID), c, nil), nil
}

func (c *channel) StageSnoop(key *ari.Key, snoopID string, opts *ari.SnoopOptions) (*ari.ChannelHandle, error) {
	k, err := c.c.getRequest(&proxy.Request{
		Kind: "ChannelStageSnoop",
		Key:  key,
		ChannelSnoop: &proxy.ChannelSnoop{
			SnoopID: snoopID,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewChannelHandle(k, c, func(h *ari.ChannelHandle) error {
		_, err := c.Snoop(k.New(ari.ChannelKey, key.ID), snoopID, opts)
		return err
	}), nil
}

func (c *channel) Dial(key *ari.Key, caller string, timeout time.Duration) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelDial",
		Key:  key,
		ChannelDial: &proxy.ChannelDial{
			Caller:  caller,
			Timeout: timeout,
		},
	})
}

func (c *channel) Play(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	k, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelPlay",
		Key:  key,
		ChannelPlay: &proxy.ChannelPlay{
			PlaybackID: playbackID,
			MediaURI:   mediaURI,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewPlaybackHandle(k.New(ari.PlaybackKey, playbackID), c.c.Playback(), nil), nil
}

func (c *channel) StagePlay(key *ari.Key, playbackID string, mediaURI string) (*ari.PlaybackHandle, error) {
	if playbackID == "" {
		playbackID = uuid.NewV1().String()
	}

	k, err := c.c.getRequest(&proxy.Request{
		Kind: "ChannelStagePlay",
		Key:  key,
		ChannelPlay: &proxy.ChannelPlay{
			PlaybackID: playbackID,
			MediaURI:   mediaURI,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewPlaybackHandle(k.New(ari.PlaybackKey, playbackID), c.c.Playback(), func(h *ari.PlaybackHandle) error {
		_, err := c.Play(k.New(ari.ChannelKey, key.ID), playbackID, mediaURI)
		return err
	}), nil
}

func (c *channel) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	rb, err := c.c.createRequest(&proxy.Request{
		Kind: "ChannelRecord",
		Key:  key,
		ChannelRecord: &proxy.ChannelRecord{
			Name:    name,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewLiveRecordingHandle(rb.New(ari.LiveRecordingKey, name), c.c.LiveRecording(), nil), nil
}

func (c *channel) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	k, err := c.c.getRequest(&proxy.Request{
		Kind: "ChannelStageRecord",
		Key:  key,
		ChannelRecord: &proxy.ChannelRecord{
			Name:    name,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewLiveRecordingHandle(k.New(ari.LiveRecordingKey, k.ID), c.c.LiveRecording(), func(h *ari.LiveRecordingHandle) error {
		_, err := c.Record(k.New(ari.ChannelKey, key.ID), k.ID, opts)
		return err
	}), nil
}

func (c *channel) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	err := c.c.commandRequest(&proxy.Request{
		Kind: "ChannelSubscribe",
		Key:  key,
	})
	if err != nil {
		c.c.log.Warn("failed to call channel subscribe")
		if key.Dialog != "" {
			c.c.log.Error("dialog present; failing")
			return nil
		}
	}
	return c.c.Bus().Subscribe(key, n...)
}

func (c *channel) GetVariable(key *ari.Key, name string) (string, error) {
	data, err := c.c.dataRequest(&proxy.Request{
		Kind: "ChannelVariableGet",
		Key:  key,
		ChannelVariable: &proxy.ChannelVariable{
			Name: name,
		},
	})
	if err != nil {
		return "", err
	}
	return data.Variable, nil
}

func (c *channel) SetVariable(key *ari.Key, name, value string) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "ChannelVariableSet",
		Key:  key,
		ChannelVariable: &proxy.ChannelVariable{
			Name:  name,
			Value: value,
		},
	})
}
