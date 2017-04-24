package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
	uuid "github.com/satori/go.uuid"
)

type bridge struct {
	c *Client
}

func (b *bridge) Create(key *ari.Key, btype, name string) (*ari.BridgeHandle, error) {
	k, err := b.c.createRequest(&proxy.Request{
		Kind: "BridgeCreate",
		Key:  key,
		BridgeCreate: &proxy.BridgeCreate{
			Type: btype,
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewBridgeHandle(k, b, nil), nil
}

func (b *bridge) StageCreate(key *ari.Key, btype, name string) (*ari.BridgeHandle, error) {
	k, err := b.c.createRequest(&proxy.Request{
		Kind: "BridgeStageCreate",
		Key:  key,
		BridgeCreate: &proxy.BridgeCreate{
			Type: btype,
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewBridgeHandle(k, b, func(h *ari.BridgeHandle) error {
		_, err := b.Create(k, btype, name)
		return err
	}), nil
}

func (b *bridge) Get(key *ari.Key) *ari.BridgeHandle {
	k, err := b.c.getRequest(&proxy.Request{
		Kind: "BridgeGet",
		Key:  key,
	})
	if err != nil {
		b.c.log.Warn("failed to get bridge for handle", "error", err)
		return ari.NewBridgeHandle(key, b, nil)
	}
	return ari.NewBridgeHandle(k, b, nil)
}

func (b *bridge) List(filter *ari.Key) ([]*ari.Key, error) {
	return b.c.listRequest(&proxy.Request{
		Kind: "BridgeList",
		Key:  filter,
	})
}

func (b *bridge) Data(key *ari.Key) (*ari.BridgeData, error) {
	resp, err := b.c.dataRequest(&proxy.Request{
		Kind: "BridgeData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Bridge, nil
}

func (b *bridge) AddChannel(key *ari.Key, channelID string) error {
	return b.c.commandRequest(&proxy.Request{
		Kind: "BridgeAddChannel",
		Key:  key,
		BridgeAddChannel: &proxy.BridgeAddChannel{
			Channel: channelID,
		},
	})
}

func (b *bridge) RemoveChannel(key *ari.Key, channelID string) error {
	return b.c.commandRequest(&proxy.Request{
		Kind: "BridgeRemoveChannel",
		Key:  key,
		BridgeRemoveChannel: &proxy.BridgeRemoveChannel{
			Channel: channelID,
		},
	})
}

func (b *bridge) Delete(key *ari.Key) error {
	return b.c.commandRequest(&proxy.Request{
		Kind: "BridgeDelete",
		Key:  key,
	})
}

func (b *bridge) Play(key *ari.Key, id string, uri string) (*ari.PlaybackHandle, error) {
	pb, err := b.c.createRequest(&proxy.Request{
		Kind: "BridgePlay",
		Key:  key,
		BridgePlay: &proxy.BridgePlay{
			MediaURI:   uri,
			PlaybackID: id,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewPlaybackHandle(pb, b.c.Playback(), nil), nil
}

func (b *bridge) StagePlay(key *ari.Key, id string, uri string) (*ari.PlaybackHandle, error) {
	k, err := b.c.getRequest(&proxy.Request{
		Kind: "BridgeStagePlay",
		Key:  key,
		BridgePlay: &proxy.BridgePlay{
			MediaURI:   uri,
			PlaybackID: id,
		},
	})
	if err != nil {
		return nil, err
	}

	return ari.NewPlaybackHandle(k, b.c.Playback(), func(h *ari.PlaybackHandle) error {
		_, err := b.Play(k, id, uri)
		return err
	}), nil
}

func (b *bridge) Record(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	if opts == nil {
		opts = &ari.RecordingOptions{}
	}
	if name == "" {
		name = uuid.NewV1().String()
	}

	rh, err := b.c.createRequest(&proxy.Request{
		Kind: "BridgeRecord",
		Key:  key,
		BridgeRecord: &proxy.BridgeRecord{
			Name:    name,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewLiveRecordingHandle(rh, b.c.LiveRecording(), nil), nil
}

func (b *bridge) StageRecord(key *ari.Key, name string, opts *ari.RecordingOptions) (*ari.LiveRecordingHandle, error) {
	if opts == nil {
		opts = &ari.RecordingOptions{}
	}
	if name == "" {
		name = uuid.NewV1().String()
	}

	k, err := b.c.getRequest(&proxy.Request{
		Kind: "BridgeStageRecord",
		Key:  key,
		BridgeRecord: &proxy.BridgeRecord{
			Name:    name,
			Options: opts,
		},
	})
	if err != nil {
		return nil, err
	}

	return ari.NewLiveRecordingHandle(k, b.c.LiveRecording(), func(h *ari.LiveRecordingHandle) error {
		_, err := b.Record(k, name, opts)
		return err
	}), nil
}

func (b *bridge) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	err := b.c.commandRequest(&proxy.Request{
		Kind: "BridgeSubscribe",
		Key:  key,
	})
	if err != nil {
		b.c.log.Warn("failed to call bridge subscribe")
		if key.Dialog != "" {
			b.c.log.Error("dialog present; failing")
			return nil
		}
	}

	return b.c.Bus().Subscribe(key, n...)
}
