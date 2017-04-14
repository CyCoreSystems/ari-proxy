package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type liveRecording struct {
	c *Client
}

func (lr *liveRecording) Get(name string) ari.LiveRecordingHandle {
	return &liveRecordingHandle{
		name: name,
		lr:   lr,
	}
}

func (lr *liveRecording) Data(name string) (lrd *ari.LiveRecordingData, err error) {
	nd, err := lr.c.dataRequest(&proxy.Request{
		RecordingLiveData: &proxy.RecordingLiveData{
			ID: name,
		},
	})
	if err != nil {
		return
	}
	lrd = nd.LiveRecording
	return
}

func (lr *liveRecording) Stop(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveStop: &proxy.RecordingLiveStop{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Pause(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLivePause: &proxy.RecordingLivePause{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Resume(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveResume: &proxy.RecordingLiveResume{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Mute(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveMute: &proxy.RecordingLiveMute{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Unmute(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveUnmute: &proxy.RecordingLiveUnmute{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Delete(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveDelete: &proxy.RecordingLiveDelete{
			ID: name,
		},
	})
	return
}

func (lr *liveRecording) Scrap(name string) (err error) {
	err = lr.c.commandRequest(&proxy.Request{
		RecordingLiveScrap: &proxy.RecordingLiveScrap{
			ID: name,
		},
	})
	return
}

type liveRecordingHandle struct {
	name string
	lr   *liveRecording
}

func (lr *liveRecordingHandle) ID() string {
	return lr.name
}

func (lr *liveRecordingHandle) Match(evt ari.Event) (ok bool) {
	return
}

func (lr *liveRecordingHandle) Data() (lrd *ari.LiveRecordingData, err error) {
	lrd, err = lr.lr.Data(lr.name)
	return
}

func (lr *liveRecordingHandle) Stop() (err error) {
	err = lr.lr.Stop(lr.name)
	return
}

func (lr *liveRecordingHandle) Pause() (err error) {
	err = lr.lr.Pause(lr.name)
	return
}

func (lr *liveRecordingHandle) Resume() (err error) {
	err = lr.lr.Resume(lr.name)
	return
}

func (lr *liveRecordingHandle) Mute() (err error) {
	err = lr.lr.Mute(lr.name)
	return
}

func (lr *liveRecordingHandle) Unmute() (err error) {
	err = lr.lr.Unmute(lr.name)
	return
}

func (lr *liveRecordingHandle) Delete() (err error) {
	err = lr.lr.Delete(lr.name)
	return
}

func (lr *liveRecordingHandle) Scrap() (err error) {
	err = lr.lr.Scrap(lr.name)
	return
}
