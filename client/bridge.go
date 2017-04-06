package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type bridge struct {
	c *Client
}

func (b *bridge) Create(id string, t string, name string) (h ari.BridgeHandle, err error) {
	req := proxy.Request{
		BridgeCreate: &proxy.BridgeCreate{
			ID:   id,
			Name: name,
			Type: t,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.GetSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = b.Get(resp.Entity.ID)
	return
}

func (b *bridge) Get(id string) ari.BridgeHandle {
	return &bridgeHandle{id: id, bridge: b}
}

func (b *bridge) List() (ret []ari.BridgeHandle, err error) {
	req := proxy.Request{
		BridgeList: &proxy.BridgeList{},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.GetSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		ret = append(ret, b.Get(i.ID))
	}
	return
}

func (b *bridge) Playback() ari.Playback {
	return b.c.Playback()
}

func (b *bridge) Data(id string) (d *ari.BridgeData, err error) {
	req := proxy.Request{
		BridgeData: &proxy.BridgeData{
			ID: id,
		},
	}
	var resp proxy.DataResponse
	err = b.c.nc.Request(proxy.GetSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	d = resp.BridgeData
	return
}

func (b *bridge) AddChannel(bridgeID string, channelID string) (err error) {
	req := proxy.Request{
		BridgeAddChannel: &proxy.BridgeAddChannel{
			ID:      bridgeID,
			Channel: channelID,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.CommandSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (b *bridge) RemoveChannel(bridgeID string, channelID string) (err error) {
	req := proxy.Request{
		BridgeRemoveChannel: &proxy.BridgeRemoveChannel{
			ID:      bridgeID,
			Channel: channelID,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.CommandSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (b *bridge) Delete(id string) (err error) {
	req := proxy.Request{
		BridgeDelete: &proxy.BridgeDelete{
			ID: id,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.CommandSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (b *bridge) Play(id string, playbackID string, mediaURI string) (h ari.PlaybackHandle, err error) {
	req := proxy.Request{
		BridgePlay: &proxy.BridgePlay{
			ID:         id,
			MediaURI:   mediaURI,
			PlaybackID: playbackID,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.CommandSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = b.c.Playback().Get(resp.Entity.ID)
	return
}

func (b *bridge) Record(id string, name string, opts *ari.RecordingOptions) (h ari.LiveRecordingHandle, err error) {
	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	req := proxy.Request{
		BridgeRecord: &proxy.BridgeRecord{
			ID:      id,
			Options: opts,
			Name:    name,
		},
	}
	var resp proxy.Response
	err = b.c.nc.Request(proxy.CommandSubject(b.c.prefix, b.c.appName, ""), &req, &resp, b.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = b.c.LiveRecording().Get(resp.Entity.ID)
	return
}

func (b *bridge) Subscribe(id string, nx ...string) ari.Subscription {
	//	ns := newSubscription(b.Get(id))
	//	ns.Start(b.subscriber, nx...)
	//	return ns
	return nil
}

type bridgeHandle struct {
	id     string
	bridge ari.Bridge
}

func (bh *bridgeHandle) ID() string {
	return bh.id
}

func (bh *bridgeHandle) Subscribe(nx ...string) ari.Subscription {
	return bh.bridge.Subscribe(bh.id, nx...)
}

func (bh *bridgeHandle) AddChannel(channelID string) error {
	return bh.bridge.AddChannel(bh.id, channelID)
}

func (bh *bridgeHandle) RemoveChannel(channelID string) error {
	return bh.bridge.AddChannel(bh.id, channelID)
}

func (bh *bridgeHandle) Delete() error {
	return bh.bridge.Delete(bh.id)
}

func (bh *bridgeHandle) Data() (bd *ari.BridgeData, err error) {
	bd, err = bh.bridge.Data(bh.id)
	return
}

func (bh *bridgeHandle) Play(playbackID string, mediaURI string) (ph ari.PlaybackHandle, err error) {
	ph, err = bh.bridge.Play(bh.id, playbackID, mediaURI)
	return
}

func (bh *bridgeHandle) Record(name string, opts *ari.RecordingOptions) (rh ari.LiveRecordingHandle, err error) {
	rh, err = bh.bridge.Record(bh.id, name, opts)
	return
}

func (bh *bridgeHandle) Match(e ari.Event) (ok bool) {
	return
}
