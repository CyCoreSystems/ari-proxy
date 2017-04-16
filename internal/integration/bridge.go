package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
	"github.com/pkg/errors"
)

func TestBridgeCreate(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var bh mocks.BridgeHandle
		bh.On("ID").Return("1234")

		m.Bridge.On("Create", "bridgeID", "bridgeType", "bridgeName").Return(&bh, nil)

		ret, err := cl.Bridge().Create("bridgeID", "bridgeType", "bridgeName")
		if err != nil {
			t.Errorf("Unexpected error in remote create call: %v", err)
		}
		if ret == nil {
			t.Errorf("Unexpected nil bridge handle")
		}
		if ret == nil || ret.ID() != bh.ID() {
			t.Errorf("Expected bridge id %v, got %v", bh.ID(), ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Create", "bridgeID", "bridgeType", "bridgeName")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("unknown error")

		m.Bridge.On("Create", "bridgeID", "bridgeType", "bridgeName").Return(nil, expected)

		ret, err := cl.Bridge().Create("bridgeID", "bridgeType", "bridgeName")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != nil {
			t.Errorf("Expected nil bridge handle, got %s", ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Create", "bridgeID", "bridgeType", "bridgeName")
	})
}

func TestBridgeList(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var handles []ari.BridgeHandle
		var h1 = &mocks.BridgeHandle{}
		h1.On("ID").Return("h1")

		var h2 = &mocks.BridgeHandle{}
		h2.On("ID").Return("h2")

		handles = append(handles, h1)
		handles = append(handles, h2)

		m.Bridge.On("List").Return(handles, nil)

		ret, err := cl.Bridge().List()
		if err != nil {
			t.Errorf("Unexpected error in remote create call: %v", err)
		}
		if len(ret) != len(handles) {
			t.Errorf("Expected handle list of length %d, got %d", len(handles), len(ret))
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "List")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("unknown error")

		m.Bridge.On("List").Return([]ari.BridgeHandle{}, expected)

		ret, err := cl.Bridge().List()
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if len(ret) != 0 {
			t.Errorf("Expected handle list of length %d, got %d", 0, len(ret))
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "List")
	})
}

func TestBridgeData(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		bd := &ari.BridgeData{
			ID:         "bridge1",
			Class:      "class1",
			ChannelIDs: []string{"channel1", "channel2"},
		}
		m.Bridge.On("Data", "1").Return(bd, nil)

		ret, err := cl.Bridge().Data("1")
		if err != nil {
			t.Errorf("Unexpected error in remote data call: %v", err)
		}
		if ret == nil || ret.ID != bd.ID || ret.Class != bd.Class {
			t.Errorf("bridge data mismatchde")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Data", "1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Data", "1").Return(nil, errors.New("Error getting data"))

		ret, err := cl.Bridge().Data("1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected bridge data to be nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Data", "1")
	})
}

func TestBridgeAddChannel(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("AddChannel", "bridge1", "channel1").Return(nil)

		err := cl.Bridge().AddChannel("bridge1", "channel1")
		if err != nil {
			t.Errorf("Unexpected error in remote AddChannel call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "AddChannel", "bridge1", "channel1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("AddChannel", "bridge1", "channel1").Return(errors.New("unknown error"))

		err := cl.Bridge().AddChannel("bridge1", "channel1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "AddChannel", "bridge1", "channel1")
	})
}

func TestBridgeRemoveChannel(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("RemoveChannel", "bridge1", "channel1").Return(nil)

		err := cl.Bridge().RemoveChannel("bridge1", "channel1")
		if err != nil {
			t.Errorf("Unexpected error in remote RemoveChannel call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "RemoveChannel", "bridge1", "channel1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("RemoveChannel", "bridge1", "channel1").Return(errors.New("unknown error"))

		err := cl.Bridge().RemoveChannel("bridge1", "channel1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "RemoveChannel", "bridge1", "channel1")
	})
}

func TestBridgeDelete(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Bridge.On("Delete", "bridge1").Return(nil)

		err := cl.Bridge().Delete("bridge1")
		if err != nil {
			t.Errorf("Unexpected error in remote Delete call: %v", err)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Delete", "bridge1")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Delete", "bridge1").Return(errors.New("unknown error"))

		err := cl.Bridge().Delete("bridge1")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Delete", "bridge1")
	})
}

func TestBridgePlay(t *testing.T, s Server) {
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var ph = &mocks.PlaybackHandle{}
		ph.On("ID").Return("playback1")

		m.Bridge.On("Play", "bridge1", "playback1", "mediaURI").Return(ph, nil)

		ret, err := cl.Bridge().Play("bridge1", "playback1", "mediaURI")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call: %v", err)
		}
		if ret == nil || ret.ID() != ph.ID() {
			t.Errorf("Expected playback handle '%v', got '%v'", ph.ID(), ret.ID())
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Play", "bridge1", "playback1", "mediaURI")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Bridge.On("Play", "bridge1", "playback1", "mediaURI").Return(nil, errors.New("unknown error"))

		ret, err := cl.Bridge().Play("bridge1", "playback1", "mediaURI")
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected empty playback handle")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Play", "bridge1", "playback1", "mediaURI")
	})
}

func TestBridgeRecord(t *testing.T, s Server) {
	runTest("customOpts", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var opts = &ari.RecordingOptions{Format: "", MaxDuration: 0, MaxSilence: 0, Exists: "", Beep: false, Terminate: "#"}

		var lrh = &mocks.LiveRecordingHandle{}
		lrh.On("ID").Return("recording1")

		m.Bridge.On("Record", "bridge1", "recording1", opts).Return(lrh, nil)

		ret, err := cl.Bridge().Record("bridge1", "recording1", opts)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %v", err)
		}
		if ret == nil || ret.ID() != lrh.ID() {
			t.Errorf("Expected liverecording handle '%v', got '%v'", lrh, ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", "bridge1", "recording1", opts)
	})

	runTest("nilOpts", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var opts = &ari.RecordingOptions{}

		var lrh = &mocks.LiveRecordingHandle{}
		lrh.On("ID").Return("recording1")

		m.Bridge.On("Record", "bridge1", "recording1", opts).Return(lrh, nil)

		ret, err := cl.Bridge().Record("bridge1", "recording1", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %v", err)
		}
		if ret == nil || ret.ID() != lrh.ID() {
			t.Errorf("Expected liverecording handle '%v', got '%v'", lrh, ret)
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", "bridge1", "recording1", opts)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var opts = &ari.RecordingOptions{}
		m.Bridge.On("Record", "bridge1", "recording1", opts).Return(nil, errors.New("unknown error"))

		ret, err := cl.Bridge().Record("bridge1", "recording1", opts)
		if err == nil {
			t.Errorf("Expected error to be non-nil")
		}
		if ret != nil {
			t.Errorf("Expected empty liverecording handle")
		}

		m.Shutdown()

		m.Bridge.AssertCalled(t, "Record", "bridge1", "recording1", opts)
	})
}
