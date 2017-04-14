package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestSoundData(t *testing.T) {
	integration.TestSoundData(t, &srv{}, clientFactory)
}

func TestSoundList(t *testing.T) {
	integration.TestSoundList(t, &srv{}, clientFactory)
}
