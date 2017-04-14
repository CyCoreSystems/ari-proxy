package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestLiveRecordingData(t *testing.T) {
	integration.TestLiveRecordingData(t, &srv{}, clientFactory)
}

func TestLiveRecordingDelete(t *testing.T) {
	integration.TestLiveRecordingDelete(t, &srv{}, clientFactory)
}

func TestLiveRecordingMute(t *testing.T) {
	integration.TestLiveRecordingMute(t, &srv{}, clientFactory)
}

func TestLiveRecordingUnmute(t *testing.T) {
	integration.TestLiveRecordingUnmute(t, &srv{}, clientFactory)
}

func TestLiveRecordingPause(t *testing.T) {
	integration.TestLiveRecordingPause(t, &srv{}, clientFactory)
}

func TestLiveRecordingStop(t *testing.T) {
	integration.TestLiveRecordingStop(t, &srv{}, clientFactory)
}

func TestLiveRecordingResume(t *testing.T) {
	integration.TestLiveRecordingResume(t, &srv{}, clientFactory)
}

func TestLiveRecordingScrap(t *testing.T) {
	integration.TestLiveRecordingScrap(t, &srv{}, clientFactory)
}
