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
	req := proxy.Request{
		PlaybackData: &proxy.PlaybackData{
			ID: id,
		},
	}
	var resp proxy.DataResponse
	err = p.c.nc.Request(proxy.GetSubject(p.c.prefix, p.c.appName, ""), &req, &resp, p.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	d = resp.PlaybackData
	return
}

func (p *playback) Control(id string, op string) (err error) {
	req := proxy.Request{
		PlaybackControl: &proxy.PlaybackControl{
			ID:      id,
			Command: op,
		},
	}
	var resp proxy.Response
	err = p.c.nc.Request(proxy.CommandSubject(p.c.prefix, p.c.appName, ""), &req, &resp, p.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (p *playback) Stop(id string) (err error) {
	req := proxy.Request{
		PlaybackStop: &proxy.PlaybackStop{
			ID: id,
		},
	}
	var resp proxy.Response
	err = p.c.nc.Request(proxy.CommandSubject(p.c.prefix, p.c.appName, ""), &req, &resp, p.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
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
