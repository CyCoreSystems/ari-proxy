package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type playback struct {
	c *Client
}

func (p *playback) Get(id string) ari.PlaybackHandle {
	return &playbackHandle{
		id:       id,
		playback: p,
	}
}

func (p *playback) Data(id string) (d *ari.PlaybackData, err error) {
	data, err := p.c.dataRequest(&proxy.Request{
		PlaybackData: &proxy.PlaybackData{
			ID: id,
		},
	})
	if err != nil {
		return
	}
	d = data.Playback
	return
}

func (p *playback) Control(id string, op string) (err error) {
	err = p.c.commandRequest(&proxy.Request{
		PlaybackControl: &proxy.PlaybackControl{
			ID:      id,
			Command: op,
		},
	})
	return
}

func (p *playback) Stop(id string) (err error) {
	err = p.c.commandRequest(&proxy.Request{
		PlaybackStop: &proxy.PlaybackStop{
			ID: id,
		},
	})
	return
}

func (p *playback) Subscribe(id string, nx ...string) ari.Subscription {
	//ns := newSubscription(p.Get(id))
	//ns.Start(p.subscriber, nx...)
	//return ns
	return nil
}

type playbackHandle struct {
	playback *playback
	id       string
}

func (ph *playbackHandle) ID() string {
	return ph.id
}

func (ph *playbackHandle) Control(op string) (err error) {
	err = ph.playback.Control(ph.id, op)
	return
}

func (ph *playbackHandle) Stop() (err error) {
	err = ph.playback.Stop(ph.id)
	return
}

func (ph *playbackHandle) Subscribe(nx ...string) ari.Subscription {
	return ph.playback.Subscribe(ph.id, nx...)
}

func (ph *playbackHandle) Data() (d *ari.PlaybackData, err error) {
	d, err = ph.playback.Data(ph.id)
	return
}

func (ph *playbackHandle) Match(evt ari.Event) (ok bool) {
	return
}
