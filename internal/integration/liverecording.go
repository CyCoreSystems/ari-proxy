package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
)

func TestLiveRecordingData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := ari.LiveRecordingData{
			Name:      "n1",
			Format:    "format",
			Cause:     "c1",
			Silence:   ari.DurationSec(3 * time.Second),
			State:     "st1",
			Talking:   ari.DurationSec(3 * time.Second),
			TargetURI: "uri1",
			Duration:  ari.DurationSec(6 * time.Second),
		}

		m.LiveRecording.On("Data", "lr1").Return(&expected, nil)

		ret, err := cl.LiveRecording().Data("lr1")
		if err != nil {
			t.Errorf("Unexpected error in liverecording Data: %s", err)
		}
		if ret == nil {
			t.Errorf("Expected live recording data to be non-nil")
		} else {
			failed := false
			failed = failed || ret.Name != expected.Name
			failed = failed || ret.Format != expected.Format
			failed = failed || ret.Cause != expected.Cause
			failed = failed || ret.Silence != expected.Silence
			failed = failed || ret.State != expected.State
			failed = failed || ret.Talking != expected.Talking
			failed = failed || ret.TargetURI != expected.TargetURI
			failed = failed || ret.Duration != expected.Duration
			if failed {
				t.Errorf("Expected '%v', got '%v'", expected, ret)
			}
		}

		m.LiveRecording.AssertCalled(t, "Data", "lr1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.LiveRecording.On("Data", "lr1").Return(nil, errors.New("err"))

		ret, err := cl.LiveRecording().Data("lr1")
		if err == nil {
			t.Errorf("Expected error in liverecording Data")
		}
		if ret != nil {
			t.Errorf("Expected live recording data to be nil")
		}

		m.LiveRecording.AssertCalled(t, "Data", "lr1")
	})
}

func testLiveRecordingCommand(t *testing.T, m *mock, name string, id string, expected error, fn func(string) error) {
	m.LiveRecording.On(name, id).Return(expected)
	err := fn(id)
	failed := false
	failed = failed || err == nil && expected != nil
	failed = failed || err != nil && expected == nil
	failed = failed || err != nil && expected != nil && err.Error() != expected.Error()
	if failed {
		t.Errorf("Expected live recording %s(%s) to return '%v', got '%v'",
			name, id, expected, err,
		)
	}
	m.LiveRecording.AssertCalled(t, name, id)
}

func TestLiveRecordingDelete(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Delete", "lr1", nil, cl.LiveRecording().Delete)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Delete", "lr1", errors.New("err"), cl.LiveRecording().Delete)
	})
}

func TestLiveRecordingMute(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Mute", "lr1", nil, cl.LiveRecording().Mute)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Mute", "lr1", errors.New("err"), cl.LiveRecording().Mute)
	})
}

func TestLiveRecordingPause(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Pause", "lr1", nil, cl.LiveRecording().Pause)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Pause", "lr1", errors.New("err"), cl.LiveRecording().Pause)
	})
}

func TestLiveRecordingStop(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Stop", "lr1", nil, cl.LiveRecording().Stop)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Stop", "lr1", errors.New("err"), cl.LiveRecording().Stop)
	})
}

func TestLiveRecordingUnmute(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Unmute", "lr1", nil, cl.LiveRecording().Unmute)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Unmute", "lr1", errors.New("err"), cl.LiveRecording().Unmute)
	})
}

func TestLiveRecordingResume(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Resume", "lr1", nil, cl.LiveRecording().Resume)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Resume", "lr1", errors.New("err"), cl.LiveRecording().Resume)
	})
}

func TestLiveRecordingScrap(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Scrap", "lr1", nil, cl.LiveRecording().Scrap)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Scrap", "lr1", errors.New("err"), cl.LiveRecording().Scrap)
	})
}
