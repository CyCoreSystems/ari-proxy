package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari/v5"
)

func TestPlaybackData(t *testing.T, s Server) {
	key := ari.NewKey(ari.PlaybackKey, "pb1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var pb ari.PlaybackData
		pb.ID = "pb1"
		pb.State = "st1"

		m.Playback.On("Data", key).Return(&pb, nil)

		ret, err := cl.Playback().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in Playback Data: %s", err)
		}
		if ret == nil {
			t.Errorf("Expected Playback data to be non-nil")
		} else {
			if ret.ID != "pb1" && ret.State != "st1" {
				t.Errorf("got '%v', expected '%v'", pb, ret)
			}
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Data", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Playback.On("Data", key).Return(nil, errors.New("error"))

		_, err := cl.Playback().Data(key)
		if err == nil {
			t.Errorf("Expected error in Playback Data: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Data", key)
	})
}

func TestPlaybackControl(t *testing.T, s Server) {
	key := ari.NewKey(ari.PlaybackKey, "pb1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Playback.On("Control", key, "op").Return(nil)

		err := cl.Playback().Control(key, "op")
		if err != nil {
			t.Errorf("Unexpected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", key, "op")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Playback.On("Control", key, "op").Return(errors.New("error"))

		err := cl.Playback().Control(key, "op")
		if err == nil {
			t.Errorf("Expected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", key, "op")
	})
}

func TestPlaybackStop(t *testing.T, s Server) {
	key := ari.NewKey(ari.PlaybackKey, "pb1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Playback.On("Stop", key).Return(nil)

		err := cl.Playback().Stop(key)
		if err != nil {
			t.Errorf("Unexpected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Playback.On("Stop", key).Return(errors.New("error"))

		err := cl.Playback().Stop(key)
		if err == nil {
			t.Errorf("Expected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", key)
	})
}
