package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestPlaybackData(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Data", "m1").Return(nil, nil)

		_, err := cl.Playback().Data("m1")
		if err != nil {
			t.Errorf("Unexpected error in Playback Data: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Data", "m1")
	})
	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Data", "m1").Return(nil, errors.New("error"))

		_, err := cl.Playback().Data("m1")
		if err == nil {
			t.Errorf("Expected error in Playback Data: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Data", "m1")
	})
}

func TestPlaybackControl(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Control", "m1", "op").Return(nil)

		err := cl.Playback().Control("m1", "op")
		if err != nil {
			t.Errorf("Unexpected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", "m1", "op")
	})
	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Control", "m1", "op").Return(errors.New("error"))

		err := cl.Playback().Control("m1", "op")
		if err == nil {
			t.Errorf("Expected error in Playback Control: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Control", "m1", "op")
	})
}

func TestPlaybackStop(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Stop", "m1").Return(nil)

		err := cl.Playback().Stop("m1")
		if err != nil {
			t.Errorf("Unexpected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", "m1")
	})
	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Playback.On("Stop", "m1").Return(errors.New("error"))

		err := cl.Playback().Stop("m1")
		if err == nil {
			t.Errorf("Expected error in Playback Stop: %s", err)
		}

		m.Shutdown()

		m.Playback.AssertCalled(t, "Stop", "m1")
	})
}
