package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestConfigData(t *testing.T) {
	integration.TestConfigData(t, &srv{}, clientFactory)
}

func TestConfigDelete(t *testing.T) {
	integration.TestConfigDelete(t, &srv{}, clientFactory)
}

func TestConfigUpdate(t *testing.T) {
	integration.TestConfigUpdate(t, &srv{}, clientFactory)
}
