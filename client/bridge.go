package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	uuid "github.com/satori/go.uuid"
)

type bridge struct {
	c *Client
}

func (b *bridge) Create(id string, t string, name string) (h ari.BridgeHandle, err error) {
	resp, err := b.c.getRequest(&proxy.Request{
		BridgeCreate: &proxy.BridgeCreate{
			ID:   id,
			Name: name,
			Type: t,
		},
	})
	if err != nil {
		return
	}
	return b.Get(resp.ID), nil
}

func (b *bridge) Get(id string) ari.BridgeHandle {
	return &bridgeHandle{id: id, bridge: b}
}

func (b *bridge) List() (ret []ari.BridgeHandle, err error) {
	resp, err := b.c.listRequest(&proxy.Request{
		BridgeList: &proxy.BridgeList{},
	})
	if err != nil {
		return nil, err
	}
	for _, i := range resp.List {
		ret = append(ret, b.Get(i.ID))
	}
	return
}

func (b *bridge) Playback() ari.Playback {
	return b.c.Playback()
}

func (b *bridge) Data(id string) (*ari.BridgeData, error) {
	resp, err := b.c.dataRequest(&proxy.Request{
		BridgeData: &proxy.BridgeData{
			ID: id,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Bridge, nil
}

func (b *bridge) AddChannel(bridgeID string, channelID string) error {
	return b.c.commandRequest(&proxy.Request{
		BridgeAddChannel: &proxy.BridgeAddChannel{
			ID:      bridgeID,
			Channel: channelID,
		},
	})
}

func (b *bridge) RemoveChannel(bridgeID string, channelID string) error {
	return b.c.commandRequest(&proxy.Request{
		BridgeRemoveChannel: &proxy.BridgeRemoveChannel{
			ID:      bridgeID,
			Channel: channelID,
		},
	})
}

func (b *bridge) Delete(id string) error {
	return b.c.commandRequest(&proxy.Request{
		BridgeDelete: &proxy.BridgeDelete{
			ID: id,
		},
	})
}

func (b *bridge) Play(id string, playbackID string, mediaURI string) (ari.PlaybackHandle, error) {
	err := b.c.commandRequest(&proxy.Request{
		BridgePlay: &proxy.BridgePlay{
			ID:         id,
			MediaURI:   mediaURI,
			PlaybackID: playbackID,
		},
	})
	if err != nil {
		return nil, err
	}
	return b.c.Playback().Get(playbackID), nil
}

func (b *bridge) Record(id string, name string, opts *ari.RecordingOptions) (ari.LiveRecordingHandle, error) {
	if opts == nil {
		opts = &ari.RecordingOptions{}
	}
	if name == "" {
		name = uuid.NewV1().String()
	}

	err := b.c.commandRequest(&proxy.Request{
		BridgeRecord: &proxy.BridgeRecord{
			ID:      id,
			Options: opts,
			Name:    name,
		},
	})
	if err != nil {
		return nil, err
	}
	return b.c.LiveRecording().Get(name), nil
}

func (b *bridge) Subscribe(id string, nx ...string) ari.Subscription {
	// TODO
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

func (bh *bridgeHandle) Data() (*ari.BridgeData, error) {
	return bh.bridge.Data(bh.id)
}

func (bh *bridgeHandle) Play(playbackID string, mediaURI string) (ari.PlaybackHandle, error) {
	return bh.bridge.Play(bh.id, playbackID, mediaURI)
}

func (bh *bridgeHandle) Record(name string, opts *ari.RecordingOptions) (ari.LiveRecordingHandle, error) {
	return bh.bridge.Record(bh.id, name, opts)
}

func (bh *bridgeHandle) Match(e ari.Event) (ok bool) {
	// TODO
	return
}
