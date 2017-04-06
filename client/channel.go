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
	return nil
}

func (c *channel) List() (ret []ari.ChannelHandle, err error) {
	req := proxy.Request{
		ChannelList: &proxy.ChannelList{},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.GetSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		ret = append(ret, c.Get(i.ID))
	}
	return
}

func (c *channel) Originate(o ari.OriginateRequest) (h ari.ChannelHandle, err error) {
	req := proxy.Request{
		ChannelOriginate: &proxy.ChannelOriginate{
			OriginateRequest: o,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CreateSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = c.Get(resp.Entity.ID)
	return
}

func (c *channel) Create(create ari.ChannelCreateRequest) (h ari.ChannelHandle, err error) {
	req := proxy.Request{
		ChannelCreate: &proxy.ChannelCreate{
			ChannelCreateRequest: create,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CreateSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = c.Get(resp.Entity.ID)
	return
}

func (c *channel) Data(id string) (cd *ari.ChannelData, err error) {
	req := proxy.Request{
		ChannelData: &proxy.ChannelData{
			ID: id,
		},
	}
	var resp proxy.DataResponse
	err = c.c.nc.Request(proxy.GetSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	cd = resp.ChannelData
	return
}

func (c *channel) Continue(id string, context string, extension string, priority int) (err error) {
	req := proxy.Request{
		ChannelContinue: &proxy.ChannelContinue{
			Context:   context,
			Extension: extension,
			Priority:  priority,
			ID:        id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Busy(id string) (err error) {
	req := proxy.Request{
		ChannelBusy: &proxy.ChannelBusy{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Congestion(id string) (err error) {
	req := proxy.Request{
		ChannelCongestion: &proxy.ChannelCongestion{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Hangup(id string, reason string) (err error) {
	req := proxy.Request{
		ChannelHangup: &proxy.ChannelHangup{
			ID:     id,
			Reason: reason,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Answer(id string) (err error) {
	req := proxy.Request{
		ChannelAnswer: &proxy.ChannelAnswer{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Ring(id string) (err error) {
	req := proxy.Request{
		ChannelRing: &proxy.ChannelRing{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) StopRing(id string) (err error) {
	req := proxy.Request{
		ChannelStopRing: &proxy.ChannelStopRing{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) SendDTMF(id string, dtmf string, opts *ari.DTMFOptions) (err error) {
	if opts == nil {
		opts = &ari.DTMFOptions{}
	}
	req := proxy.Request{
		ChannelSendDTMF: &proxy.ChannelSendDTMF{
			ID:      id,
			DTMF:    dtmf,
			Options: opts,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Hold(id string) (err error) {
	req := proxy.Request{
		ChannelHold: &proxy.ChannelHold{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) StopHold(id string) (err error) {
	req := proxy.Request{
		ChannelStopHold: &proxy.ChannelStopHold{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Mute(id string, dir ari.Direction) (err error) {
	req := proxy.Request{
		ChannelMute: &proxy.ChannelMute{
			ID:        id,
			Direction: dir,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Unmute(id string, dir ari.Direction) (err error) {
	req := proxy.Request{
		ChannelUnmute: &proxy.ChannelUnmute{
			ID:        id,
			Direction: dir,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) MOH(id string, moh string) (err error) {
	req := proxy.Request{
		ChannelMOH: &proxy.ChannelMOH{
			ID:    id,
			Music: moh,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) StopMOH(id string) (err error) {
	req := proxy.Request{
		ChannelStopMOH: &proxy.ChannelStopMOH{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Silence(id string) (err error) {
	req := proxy.Request{
		ChannelSilence: &proxy.ChannelSilence{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) StopSilence(id string) (err error) {
	req := proxy.Request{
		ChannelStopSilence: &proxy.ChannelStopSilence{
			ID: id,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

// SnoopRequest is the NATs snoop request
type SnoopRequest struct {
	SnoopID string
	App     string
	Options *ari.SnoopOptions
}

func (c *channel) Snoop(id string, snoopID string, opts *ari.SnoopOptions) (ch ari.ChannelHandle, err error) {
	req := proxy.Request{
		ChannelSnoop: &proxy.ChannelSnoop{
			ID:      id,
			SnoopID: snoopID,
			Options: opts,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	ch = c.Get(resp.Entity.ID)
	return
}

// DialRequest is the request for the channel dial operation
type DialRequest struct {
	Caller  string `json:"caller"`
	Timeout int    `json:"timeout"`
}

func (c *channel) Dial(id string, caller string, timeout time.Duration) (err error) {
	req := proxy.Request{
		ChannelDial: &proxy.ChannelDial{
			ID:      id,
			Caller:  caller,
			Timeout: timeout,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *channel) Play(id string, playbackID string, mediaURI string) (p ari.PlaybackHandle, err error) {
	req := proxy.Request{
		ChannelPlay: &proxy.ChannelPlay{
			ID:         id,
			PlaybackID: playbackID,
			MediaURI:   mediaURI,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	p = c.Playback().Get(resp.Entity.ID)
	return
}

func (c *channel) Record(id string, name string, opts *ari.RecordingOptions) (h ari.LiveRecordingHandle, err error) {
	req := proxy.Request{
		ChannelRecord: &proxy.ChannelRecord{
			ID:      id,
			Name:    name,
			Options: opts,
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.CommandSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = c.c.LiveRecording().Get(resp.Entity.ID)
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
	req := proxy.Request{
		ChannelVariables: &proxy.ChannelVariables{
			Name: variable,
			ID:   c.id,
			Get:  &proxy.VariablesGet{},
		},
	}
	var resp proxy.DataResponse
	err = c.c.c.nc.Request(proxy.GetSubject(c.c.c.prefix, c.c.c.appName, ""), &req, &resp, c.c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	ret = resp.Variable
	return
}

func (c *channelVariables) Set(variable string, value string) (err error) {
	req := proxy.Request{
		ChannelVariables: &proxy.ChannelVariables{
			Name: variable,
			ID:   c.id,
			Set: &proxy.VariablesSet{
				Value: value,
			},
		},
	}
	var resp proxy.Response
	err = c.c.c.nc.Request(proxy.CommandSubject(c.c.c.prefix, c.c.c.appName, ""), &req, &resp, c.c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}
