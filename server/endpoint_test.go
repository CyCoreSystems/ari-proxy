package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestEndpointList(t *testing.T) {
	integration.TestEndpointList(t, &srv{}, clientFactory)
}

func TestEndpointListByTech(t *testing.T) {
	integration.TestEndpointListByTech(t, &srv{}, clientFactory)
}

func TestEndpointData(t *testing.T) {
	integration.TestEndpointData(t, &srv{}, clientFactory)
}
