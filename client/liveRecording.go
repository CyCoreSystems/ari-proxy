package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type liveRecording struct {
	c *Client
}

func (l *liveRecording) Get(key *ari.Key) *ari.LiveRecordingHandle {
	k, err := l.c.getRequest(&proxy.Request{
		Kind: "RecordingLiveGet",
		Key:  key,
	})
	if err != nil {
		l.c.log.Warn("failed to get liveRecording for handle", "error", err)
		return ari.NewLiveRecordingHandle(key, l, nil)
	}
	return ari.NewLiveRecordingHandle(k, l, nil)

}

func (l *liveRecording) Data(key *ari.Key) (*ari.LiveRecordingData, error) {
	data, err := l.c.dataRequest(&proxy.Request{
		Kind: "RecordingLiveData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.LiveRecording, nil
}

func (l *liveRecording) Stop(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveStop",
		Key:  key,
	})
}

func (l *liveRecording) Pause(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLivePause",
		Key:  key,
	})
}

func (l *liveRecording) Resume(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveResume",
		Key:  key,
	})
}

func (l *liveRecording) Mute(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveMute",
		Key:  key,
	})
}

func (l *liveRecording) Unmute(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveUnmute",
		Key:  key,
	})
}

func (l *liveRecording) Scrap(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveScrap",
		Key:  key,
	})
}

func (l *liveRecording) Stored(key *ari.Key) *ari.StoredRecordingHandle {
	return ari.NewStoredRecordingHandle(key.New(ari.StoredRecordingKey, key.ID), l.c.StoredRecording(), nil)
}

func (l *liveRecording) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	err := l.c.commandRequest(&proxy.Request{
		Kind: "RecordingLiveSubscribe",
		Key:  key,
	})
	if err != nil {
		l.c.log.Warn("failed to call liveRecording Subscribe")
		if key.Dialog != "" {
			l.c.log.Error("dialog present; failing")
			return nil
		}
	}

	return l.c.Bus().Subscribe(key, n...)
}
