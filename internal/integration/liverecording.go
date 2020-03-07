package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari/v5"
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

		key := ari.NewKey(ari.LiveRecordingKey, "lr1")

		m.LiveRecording.On("Data", key).Return(&expected, nil)

		ret, err := cl.LiveRecording().Data(key)
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

		m.LiveRecording.AssertCalled(t, "Data", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		key := ari.NewKey(ari.LiveRecordingKey, "lr1")

		m.LiveRecording.On("Data", key).Return(nil, errors.New("err"))

		ret, err := cl.LiveRecording().Data(key)
		if err == nil {
			t.Errorf("Expected error in liverecording Data")
		}
		if ret != nil {
			t.Errorf("Expected live recording data to be nil")
		}

		m.LiveRecording.AssertCalled(t, "Data", key)
	})
}

func testLiveRecordingCommand(t *testing.T, m *mock, name string, id *ari.Key, expected error, fn func(*ari.Key) error) {
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

/*
func TestLiveRecordingDelete(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Delete", key, nil, cl.LiveRecording().Delete)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Delete", key, errors.New("err"), cl.LiveRecording().Delete)
	})
}
*/

func TestLiveRecordingMute(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Mute", key, nil, cl.LiveRecording().Mute)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Mute", key, errors.New("err"), cl.LiveRecording().Mute)
	})
}

func TestLiveRecordingPause(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Pause", key, nil, cl.LiveRecording().Pause)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Pause", key, errors.New("err"), cl.LiveRecording().Pause)
	})
}

func TestLiveRecordingStop(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Stop", key, nil, cl.LiveRecording().Stop)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Stop", key, errors.New("err"), cl.LiveRecording().Stop)
	})
}

func TestLiveRecordingUnmute(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Unmute", key, nil, cl.LiveRecording().Unmute)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Unmute", key, errors.New("err"), cl.LiveRecording().Unmute)
	})
}

func TestLiveRecordingResume(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Resume", key, nil, cl.LiveRecording().Resume)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Resume", key, errors.New("err"), cl.LiveRecording().Resume)
	})
}

func TestLiveRecordingScrap(t *testing.T, s Server) {
	key := ari.NewKey(ari.LiveRecordingKey, "lr1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Scrap", key, nil, cl.LiveRecording().Scrap)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		testLiveRecordingCommand(t, m, "Scrap", key, errors.New("err"), cl.LiveRecording().Scrap)
	})
}
