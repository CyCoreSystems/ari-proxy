package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestPlaybackData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var pb ari.PlaybackData
		pb.ID = "pb1"
		pb.State = "st1"

		m.Playback.On("Data", "pb1").Return(&pb, nil)

		ret, err := cl.Playback().Data("pb1")
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

		m.Playback.AssertCalled(t, "Data", "pb1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Data", "pb1").Return(nil, errors.New("error"))

		_, err := cl.Playback().Data("pb1")
		if err == nil {
			t.Errorf("Expected error in Playback Data: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Data", "pb1")
	})
}

func TestPlaybackControl(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Control", "m1", "op").Return(nil)

		err := cl.Playback().Control("m1", "op")
		if err != nil {
			t.Errorf("Unexpected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", "m1", "op")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Control", "m1", "op").Return(errors.New("error"))

		err := cl.Playback().Control("m1", "op")
		if err == nil {
			t.Errorf("Expected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", "m1", "op")
	})
}

func TestPlaybackStop(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Stop", "m1").Return(nil)

		err := cl.Playback().Stop("m1")
		if err != nil {
			t.Errorf("Unexpected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", "m1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Stop", "m1").Return(errors.New("error"))

		err := cl.Playback().Stop("m1")
		if err == nil {
			t.Errorf("Expected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", "m1")
	})
}
