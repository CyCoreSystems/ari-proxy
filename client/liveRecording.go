package client

import "github.com/CyCoreSystems/ari"

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
	return
}

func (lr *liveRecording) Stop(name string) (err error) {
	return
}

func (lr *liveRecording) Pause(name string) (err error) {
	return
}

func (lr *liveRecording) Resume(name string) (err error) {
	return
}

func (lr *liveRecording) Mute(name string) (err error) {
	return
}

func (lr *liveRecording) Unmute(name string) (err error) {
	return
}

func (lr *liveRecording) Delete(name string) (err error) {
	return
}

func (lr *liveRecording) Scrap(name string) (err error) {
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
	return
}

func (lr *liveRecordingHandle) Stop() (err error) {
	return
}

func (lr *liveRecordingHandle) Pause() (err error) {
	return
}

func (lr *liveRecordingHandle) Resume() (err error) {
	return
}

func (lr *liveRecordingHandle) Mute() (err error) {
	return
}

func (lr *liveRecordingHandle) Unmute() (err error) {
	return
}

func (lr *liveRecordingHandle) Delete() (err error) {
	return
}

func (lr *liveRecordingHandle) Scrap() (err error) {
	return
}
