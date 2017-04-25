package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestSoundData(t *testing.T) {
	integration.TestSoundData(t, &srv{})
}

func TestSoundList(t *testing.T) {
	integration.TestSoundList(t, &srv{})
}
