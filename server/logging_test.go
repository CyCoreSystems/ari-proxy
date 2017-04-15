package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestLoggingList(t *testing.T) {
	integration.TestLoggingList(t, &srv{}, clientFactory)
}

func TestLoggingCreate(t *testing.T) {
	integration.TestLoggingCreate(t, &srv{}, clientFactory)
}

func TestLoggingRotate(t *testing.T) {
	integration.TestLoggingRotate(t, &srv{}, clientFactory)
}

func TestLoggingDelete(t *testing.T) {
	integration.TestLoggingDelete(t, &srv{}, clientFactory)
}
