package client

import (
	"time"

	"github.com/CyCoreSystems/ari"
)

type natsBridge struct {
	conn          *Conn
	subscriber    ari.Subscriber
	playback      ari.Playback
	liveRecording ari.LiveRecording
}

// CreateBridgeRequest is the request for creating bridges
type CreateBridgeRequest struct {
	ID   string `json:"bridgeId,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func (b *natsBridge) Create(id string, t string, name string) (h *ari.BridgeHandle, err error) {
	var bridgeID string
	req := CreateBridgeRequest{id, t, name}
	err = b.conn.StandardRequest("ari.bridges.create", "", &req, &bridgeID)
	if err != nil {
		return
	}
	h = b.Get(bridgeID)
	return
}

func (b *natsBridge) Get(id string) *ari.BridgeHandle {
	return ari.NewBridgeHandle(id, b)
}

func (b *natsBridge) List() (bx []*ari.BridgeHandle, err error) {
	var bridges []string
	err = b.conn.ReadRequest("ari.bridges.all", "", nil, &bridges)
	for _, bridge := range bridges {
		bx = append(bx, b.Get(bridge))
	}
	return
}

func (b *natsBridge) Playback() ari.Playback {
	return b.playback
}

func (b *natsBridge) Data(id string) (d ari.BridgeData, err error) {
	err = b.conn.ReadRequest("ari.bridges.data", id, nil, &d)
	return
}

func (b *natsBridge) AddChannel(bridgeID string, channelID string) (err error) {
	err = b.conn.StandardRequest("ari.bridges.addChannel", bridgeID, channelID, nil)
	return
}

func (b *natsBridge) RemoveChannel(bridgeID string, channelID string) (err error) {
	err = b.conn.StandardRequest("ari.bridges.removeChannel", bridgeID, channelID, nil)
	return
}

func (b *natsBridge) Delete(id string) (err error) {
	err = b.conn.StandardRequest("ari.bridges.delete", id, nil, nil)
	return
}

// PlayRequest is the request for playback
type PlayRequest struct {
	PlaybackID string `json:"playback_id"`
	MediaURI   string `json:"media_uri"`
}

func (b *natsBridge) Play(id string, playbackID string, mediaURI string) (h *ari.PlaybackHandle, err error) {
	err = b.conn.StandardRequest("ari.bridges.play", id, &PlayRequest{PlaybackID: playbackID, MediaURI: mediaURI}, nil)
	if err == nil {
		h = b.playback.Get(playbackID)
	}
	return
}

// RecordRequest is a request for recording
type RecordRequest struct {
	Name        string `json:"name"`
	Format      string `json:"format"`
	MaxDuration int    `json:"maxDurationSeconds"`
	MaxSilence  int    `json:"maxSilenceSeconds"`
	IfExists    string `json:"ifExists,omitempty"`
	Beep        bool   `json:"beep"`
	TerminateOn string `json:"terminateOn,omitempty"`
}

func (b *natsBridge) Record(id string, name string, opts *ari.RecordingOptions) (h *ari.LiveRecordingHandle, err error) {

	if opts == nil {
		opts = &ari.RecordingOptions{}
	}

	req := RecordRequest{
		Name:        name,
		Format:      opts.Format,
		MaxDuration: int(opts.MaxDuration / time.Second),
		MaxSilence:  int(opts.MaxSilence / time.Second),
		IfExists:    opts.Exists,
		Beep:        opts.Beep,
		TerminateOn: opts.Terminate,
	}
	err = b.conn.StandardRequest("ari.bridges.record", id, req, nil)
	if err == nil {
		h = b.liveRecording.Get(name)
	}
	return
}

func (b *natsBridge) Subscribe(id string, nx ...string) ari.Subscription {
	ns := newSubscription(b.Get(id))
	go ns.Run(b.subscriber, nx...)
	return ns
}
