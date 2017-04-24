package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/pkg/errors"
)

func TestBridgeCreate(t *testing.T, s Server) {

	key := ari.NewKey(ari.BridgeKey, "bridgeID")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		bh := ari.NewBridgeHandle(key, m.Bridge, nil)

		m.Bridge.On("Create", key, "bridgeType", "bridgeName").Return(bh, nil)

		ret, err := cl.Bridge().Create(key, "bridgeType", "bridgeName")
		if err != nil {
			t.Errorf("Unexpected error in remote create call: %v", err)
		}
		if ret == nil {
			t.Errorf("Unexpected nil bridge handle")
		}
		if ret == nil || ret.ID() != key.ID {
			t.Errorf("Expected bridge id %v, got %v", key.ID, ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Create", key, "bridgeType", "bridgeName")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("unknown error")

		m.Bridge.On("Create", key, "bridgeType", "bridgeName").Return(nil, expected)

		ret, err := cl.Bridge().Create(key, "bridgeType", "bridgeName")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != nil {
			t.Errorf("Expected nil bridge handle, got %v", ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Create", key, "bridgeType", "bridgeName")
	})
}

func TestBridgeList(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var handles []*ari.Key
		var h1 = ari.NewKey(ari.BridgeKey, "h1")
		var h2 = ari.NewKey(ari.BridgeKey, "h2")

		handles = append(handles, h1)
		handles = append(handles, h2)

		m.Bridge.On("List", (*ari.Key)(nil)).Return(handles, nil)

		ret, err := cl.Bridge().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote create call: %v", err)
		}
		if len(ret) != len(handles) {
			t.Errorf("Expected handle list of length %d, got %d", len(handles), len(ret))
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("unknown error")

		m.Bridge.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, expected)

		ret, err := cl.Bridge().List(nil)
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if len(ret) != 0 {
			t.Errorf("Expected handle list of length %d, got %d", 0, len(ret))
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "List", (*ari.Key)(nil))
	})
}

func TestBridgeData(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		bd := &ari.BridgeData{
			ID:         "bridge1",
			Class:      "class1",
			ChannelIDs: []string{"channel1", "channel2"},
		}
		m.Bridge.On("Data", key).Return(bd, nil)

		ret, err := cl.Bridge().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in remote data call: %v", err)
		}
		if ret == nil || ret.ID != bd.ID || ret.Class != bd.Class {
			t.Errorf("bridge data mismatchde")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Data", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Data", key).Return(nil, errors.New("Error getting data"))

		ret, err := cl.Bridge().Data(key)
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected bridge data to be nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Data", key)
	})
}

func TestBridgeAddChannel(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("AddChannel", key, "channel1").Return(nil)

		err := cl.Bridge().AddChannel(key, "channel1")
		if err != nil {
			t.Errorf("Unexpected error in remote AddChannel call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "AddChannel", key, "channel1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("AddChannel", key, "channel1").Return(errors.New("unknown error"))

		err := cl.Bridge().AddChannel(key, "channel1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "AddChannel", key, "channel1")
	})
}

func TestBridgeRemoveChannel(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("RemoveChannel", key, "channel1").Return(nil)

		err := cl.Bridge().RemoveChannel(key, "channel1")
		if err != nil {
			t.Errorf("Unexpected error in remote RemoveChannel call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "RemoveChannel", key, "channel1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("RemoveChannel", key, "channel1").Return(errors.New("unknown error"))

		err := cl.Bridge().RemoveChannel(key, "channel1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "RemoveChannel", key, "channel1")
	})
}

func TestBridgeDelete(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("Delete", key).Return(nil)

		err := cl.Bridge().Delete(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Delete call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Delete", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Delete", key).Return(errors.New("unknown error"))

		err := cl.Bridge().Delete(key)
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Delete", key)
	})
}

func TestBridgePlay(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var ph = ari.NewPlaybackHandle(ari.NewKey(ari.PlaybackKey, "playback1"), m.Playback, nil)

		m.Bridge.On("Play", key, "playback1", "mediaURI").Return(ph, nil)

		ret, err := cl.Bridge().Play(key, "playback1", "mediaURI")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call: %v", err)
		}
		if ret == nil || ret.ID() != ph.ID() {
			t.Errorf("Expected playback handle '%v', got '%v'", ph.ID(), ret.ID())
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Play", key, "playback1", "mediaURI")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Play", key, "playback1", "mediaURI").Return(nil, errors.New("unknown error"))

		ret, err := cl.Bridge().Play(key, "playback1", "mediaURI")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected empty playback handle")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Play", key, "playback1", "mediaURI")
	})
}

func TestBridgeRecord(t *testing.T, s Server) {
	var key = ari.NewKey(ari.BridgeKey, "bridge1")

	runTest("customOpts", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var opts = &ari.RecordingOptions{Format: "", MaxDuration: 0, MaxSilence: 0, Exists: "", Beep: false, Terminate: "#"}

		var lrh = ari.NewLiveRecordingHandle(ari.NewKey(ari.LiveRecordingKey, "recording1"), m.LiveRecording, nil)

		bd := &ari.BridgeData{
			ID:         "bridge1",
			Class:      "class1",
			ChannelIDs: []string{"channel1", "channel2"},
		}
		m.Bridge.On("Data", key).Return(bd, nil)
		m.Bridge.On("Record", key, "recording1", opts).Return(lrh, nil)

		ret, err := cl.Bridge().Record(key, "recording1", opts)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %v", err)
		}
		if ret == nil || ret.ID() != lrh.ID() {
			t.Errorf("Expected liverecording handle '%v', got '%v'", lrh.ID(), ret.ID())
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", key, "recording1", opts)
	})

	runTest("nilOpts", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var opts = &ari.RecordingOptions{}

		var lrh = ari.NewLiveRecordingHandle(ari.NewKey(ari.LiveRecordingKey, "recording1"), m.LiveRecording, nil)

		bd := &ari.BridgeData{
			ID:         "bridge1",
			Class:      "class1",
			ChannelIDs: []string{"channel1", "channel2"},
		}
		m.Bridge.On("Data", key).Return(bd, nil)
		m.Bridge.On("Record", key, "recording1", opts).Return(lrh, nil)

		ret, err := cl.Bridge().Record(key, "recording1", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %v", err)
		}
		if ret == nil || ret.ID() != lrh.ID() {
			t.Errorf("Expected liverecording handle '%v', got '%v'", lrh, ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", key, "recording1", opts)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		bd := &ari.BridgeData{
			ID:         "bridge1",
			Class:      "class1",
			ChannelIDs: []string{"channel1", "channel2"},
		}
		m.Bridge.On("Data", key).Return(bd, nil)

		var opts = &ari.RecordingOptions{}
		m.Bridge.On("Record", key, "recording1", opts).Return(nil, errors.New("unknown error"))

		ret, err := cl.Bridge().Record(key, "recording1", opts)
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected empty liverecording handle")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", key, "recording1", opts)
	})
}
