package client

import (
	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type playback struct {
	c *Client
}

func (p *playback) Get(key *ari.Key) *ari.PlaybackHandle {
	k, err := p.c.getRequest(&proxy.Request{
		Kind: "PlaybackGet",
		Key:  key,
	})
	if err != nil {
		p.c.log.Warn("failed to get playback for handle", "error", err)
		return ari.NewPlaybackHandle(key, p, nil)
	}
	return ari.NewPlaybackHandle(k, p, nil)
}

func (p *playback) Data(key *ari.Key) (*ari.PlaybackData, error) {
	data, err := p.c.dataRequest(&proxy.Request{
		Kind: "PlaybackData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Playback, nil
}

func (p *playback) Control(key *ari.Key, op string) error {
	return p.c.commandRequest(&proxy.Request{
		Kind: "PlaybackControl",
		Key:  key,
		PlaybackControl: &proxy.PlaybackControl{
			Command: op,
		},
	})
}

func (p *playback) Stop(key *ari.Key) error {
	return p.c.commandRequest(&proxy.Request{
		Kind: "PlaybackStop",
		Key:  key,
	})
}

func (p *playback) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	err := p.c.commandRequest(&proxy.Request{
		Kind: "PlaybackSubscribe",
		Key:  key,
	})
	if err != nil {
		p.c.log.Warn("failed to call bridge subscribe", "error", err)
		if key.Dialog != "" {
			p.c.log.Error("dialog present; failing", "error", err)
			return nil
		}
	}
	return p.c.Bus().Subscribe(key, n...)
}
