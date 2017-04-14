package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestPlaybackData(t *testing.T) {
	integration.TestPlaybackData(t, &srv{}, clientFactory)
}

func TestPlaybackControl(t *testing.T) {
	integration.TestPlaybackControl(t, &srv{}, clientFactory)
}

func TestPlaybackStop(t *testing.T) {
	integration.TestPlaybackStop(t, &srv{}, clientFactory)
}
