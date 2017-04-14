package client

import (
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type channel struct {
	c *Client
}

func (c *channel) Playback() ari.Playback {
	return c.c.Playback()
}

func (c *channel) Get(id string) ari.ChannelHandle {
	return &channelHandle{
		id:      id,
		channel: c,
	}
}

func (c *channel) List() (ret []ari.ChannelHandle, err error) {
	el, err := c.c.listRequest(&proxy.Request{
		ChannelList: &proxy.ChannelList{},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
		ret = append(ret, c.Get(i.ID))
	}
	return
}

func (c *channel) Originate(o ari.OriginateRequest) (h ari.ChannelHandle, err error) {
	e, err := c.c.createRequest(&proxy.Request{
		ChannelOriginate: &proxy.ChannelOriginate{
			OriginateRequest: o,
		},
	})
	if err != nil {
		return
	}
	h = c.Get(e.ID)
	return
}

func (c *channel) Create(create ari.ChannelCreateRequest) (h ari.ChannelHandle, err error) {
	e, err := c.c.createRequest(&proxy.Request{
		ChannelCreate: &proxy.ChannelCreate{
			ChannelCreateRequest: create,
		},
	})
	if err != nil {
		return
	}
	h = c.Get(e.ID)
	return
}

func (c *channel) Data(id string) (cd *ari.ChannelData, err error) {
	d, err := c.c.dataRequest(&proxy.Request{
		ChannelData: &proxy.ChannelData{
			ID: id,
		},
	})
	if err != nil {
		return
	}
	cd = d.Channel
	return
}

func (c *channel) Continue(id string, context string, extension string, priority int) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelContinue: &proxy.ChannelContinue{
			Context:   context,
			Extension: extension,
			Priority:  priority,
			ID:        id,
		},
	})
	return
}

func (c *channel) Busy(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelBusy: &proxy.ChannelBusy{
			ID: id,
		},
	})
	return
}

func (c *channel) Congestion(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelCongestion: &proxy.ChannelCongestion{
			ID: id,
		},
	})
	return
}

func (c *channel) Hangup(id string, reason string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelHangup: &proxy.ChannelHangup{
			ID:     id,
			Reason: reason,
		},
	})
	return
}

func (c *channel) Answer(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelAnswer: &proxy.ChannelAnswer{
			ID: id,
		},
	})
	return
}

func (c *channel) Ring(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelRing: &proxy.ChannelRing{
			ID: id,
		},
	})
	return
}

func (c *channel) StopRing(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelStopRing: &proxy.ChannelStopRing{
			ID: id,
		},
	})
	return
}

func (c *channel) SendDTMF(id string, dtmf string, opts *ari.DTMFOptions) (err error) {
	if opts == nil {
		opts = &ari.DTMFOptions{}
	}
	err = c.c.commandRequest(&proxy.Request{
		ChannelSendDTMF: &proxy.ChannelSendDTMF{
			ID:      id,
			DTMF:    dtmf,
			Options: opts,
		},
	})
	return
}

func (c *channel) Hold(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelHold: &proxy.ChannelHold{
			ID: id,
		},
	})
	return
}

func (c *channel) StopHold(id string) (err error) {
	err = c.c.commandRequest(proxy.Request{
		ChannelStopHold: &proxy.ChannelStopHold{
			ID: id,
		},
	})
	return
}

func (c *channel) Mute(id string, dir ari.Direction) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelMute: &proxy.ChannelMute{
			ID:        id,
			Direction: dir,
		},
	})
	return
}

func (c *channel) Unmute(id string, dir ari.Direction) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelUnmute: &proxy.ChannelUnmute{
			ID:        id,
			Direction: dir,
		},
	})
	return
}

func (c *channel) MOH(id string, moh string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelMOH: &proxy.ChannelMOH{
			ID:    id,
			Music: moh,
		},
	})
	return
}

func (c *channel) StopMOH(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelStopMOH: &proxy.ChannelStopMOH{
			ID: id,
		},
	})
	return
}

func (c *channel) Silence(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelSilence: &proxy.ChannelSilence{
			ID: id,
		},
	})
	return
}

func (c *channel) StopSilence(id string) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelStopSilence: &proxy.ChannelStopSilence{
			ID: id,
		},
	})
	return
}

func (c *channel) Snoop(id string, snoopID string, opts *ari.SnoopOptions) (ch ari.ChannelHandle, err error) {
	e, err := c.c.createRequest(&proxy.Request{
		ChannelSnoop: &proxy.ChannelSnoop{
			ID:      id,
			SnoopID: snoopID,
			Options: opts,
		},
	})
	if err != nil {
		return
	}
	ch = c.Get(e.ID)
	return
}

func (c *channel) Dial(id string, caller string, timeout time.Duration) (err error) {
	err = c.c.commandRequest(&proxy.Request{
		ChannelDial: &proxy.ChannelDial{
			ID:      id,
			Caller:  caller,
			Timeout: timeout,
		},
	})
	return
}

func (c *channel) Play(id string, playbackID string, mediaURI string) (p ari.PlaybackHandle, err error) {
	e, err := c.c.createRequest(&proxy.Request{
		ChannelPlay: &proxy.ChannelPlay{
			ID:         id,
			PlaybackID: playbackID,
			MediaURI:   mediaURI,
		},
	})
	if err != nil {
		return
	}
	p = c.Playback().Get(e.ID)
	return
}

func (c *channel) Record(id string, name string, opts *ari.RecordingOptions) (h ari.LiveRecordingHandle, err error) {
	e, err := c.c.createRequest(&proxy.Request{ChannelRecord: &proxy.ChannelRecord{
		ID:      id,
		Name:    name,
		Options: opts,
	},
	})
	if err != nil {
		return
	}
	h = c.c.LiveRecording().Get(e.ID)
	return
}

func (c *channel) Subscribe(id string, n ...string) ari.Subscription {
	//ns := newSubscription(c.Get(id))
	//ns.Start(c.subscriber, n...)
	//return ns
	return nil
}

type channelVariables struct {
	c  *channel
	id string
}

func (c *channel) Variables(id string) ari.Variables {
	return &channelVariables{c, id}
}

// GetChannelVariable is the request object for getting a channel variable
type GetChannelVariable struct {
	Name string `json:"name"`
}

// SetChannelVariable is the request object for setting a channel variable
type SetChannelVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *channelVariables) Get(variable string) (ret string, err error) {
	ed, err := c.c.c.dataRequest(&proxy.Request{
		ChannelVariables: &proxy.ChannelVariables{
			Name: variable,
			ID:   c.id,
			Get:  &proxy.VariablesGet{},
		},
	})
	if err != nil {
		return
	}
	ret = ed.Variable
	return
}

func (c *channelVariables) Set(variable string, value string) (err error) {
	err = c.c.c.commandRequest(&proxy.Request{
		ChannelVariables: &proxy.ChannelVariables{
			Name: variable,
			ID:   c.id,
			Set: &proxy.VariablesSet{
				Value: value,
			},
		},
	})
	return
}

type channelHandle struct {
	id      string
	channel *channel
}

func (c *channelHandle) Answer() error {
	return c.channel.Answer(c.id)
}

func (c *channelHandle) Busy() error {
	return c.channel.Busy(c.id)
}

func (c *channelHandle) Congestion() error {
	return c.channel.Congestion(c.id)
}

func (c *channelHandle) Continue(s string, s1 string, i int) error {
	return c.channel.Continue(c.id, s, s1, i)
}

func (c *channelHandle) Data() (*ari.ChannelData, error) {
	return c.channel.Data(c.id)
}

func (c *channelHandle) Dial(s string, d time.Duration) error {
	return c.channel.Dial(c.id, s, d)
}

func (c *channelHandle) Hangup() error {
	return c.channel.Hangup(c.id, "")
}

func (c *channelHandle) Hold() error {
	return c.channel.Hold(c.id)
}

func (c *channelHandle) ID() string {
	return c.id
}

func (c *channelHandle) IsAnswered() (bool, error) {
	d, err := c.Data()
	if err != nil {
		return false, err
	}
	if d.State != "Up" {
		return false, nil
	}
	return true, nil
}

func (c *channelHandle) MOH(s string) error {
	return c.channel.MOH(c.id, s)
}

func (c *channelHandle) Match(e ari.Event) bool {
	v, ok := e.(ari.ChannelEvent)
	if !ok {
		return false
	}
	list := v.GetChannelIDs()
	for _, i := range list {
		if i == c.id {
			return true
		}
	}
	return false
}

func (c *channelHandle) Mute(d ari.Direction) error {
	return c.channel.Mute(c.id, d)
}

func (c *channelHandle) Originate(o ari.OriginateRequest) (ari.ChannelHandle, error) {
	if o.ChannelID == "" {
		o.ChannelID = c.ID()
	}
	return c.channel.Originate(o)
}

func (c *channelHandle) Play(playbackID string, mediaURI string) (ari.PlaybackHandle, error) {
	return c.channel.Play(c.id, playbackID, mediaURI)
}

func (c *channelHandle) Record(name string, r *ari.RecordingOptions) (ari.LiveRecordingHandle, error) {
	return c.channel.Record(c.id, name, r)
}

func (c *channelHandle) Ring() error {
	return c.channel.Ring(c.id)
}

func (c *channelHandle) SendDTMF(s string, d *ari.DTMFOptions) error {
	return c.channel.SendDTMF(c.id, s, d)
}

func (c *channelHandle) Silence() error {
	return c.channel.Silence(c.id)
}

func (c *channelHandle) Snoop(snoopID string, opts *ari.SnoopOptions) (ari.ChannelHandle, error) {
	return c.channel.Snoop(c.id, snoopID, opts)
}

func (c *channelHandle) StopHold() error {
	return c.channel.StopHold(c.id)
}

func (c *channelHandle) StopMOH() error {
	return c.channel.StopMOH(c.id)
}

func (c *channelHandle) StopRing() error {
	return c.channel.StopRing(c.id)
}

func (c *channelHandle) StopSilence() error {
	return c.channel.StopSilence(c.id)
}

func (c *channelHandle) Subscribe(nx ...string) ari.Subscription {
	return c.channel.Subscribe(c.id, nx...)
}

func (c *channelHandle) Unmute(d ari.Direction) error {
	return c.channel.Unmute(c.id, d)
}

func (c *channelHandle) Variables() ari.Variables {
	return c.channel.Variables(c.id)
}
