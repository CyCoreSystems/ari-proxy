package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestDeviceData(t *testing.T) {
	integration.TestDeviceData(t, &srv{}, clientFactory)
}

func TestDeviceDelete(t *testing.T) {
	integration.TestDeviceDelete(t, &srv{}, clientFactory)
}

func TestDeviceUpdate(t *testing.T) {
	integration.TestDeviceUpdate(t, &srv{}, clientFactory)
}

func TestDeviceList(t *testing.T) {
	integration.TestDeviceList(t, &srv{}, clientFactory)
}
