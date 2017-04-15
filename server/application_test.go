package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestApplicationList(t *testing.T) {
	integration.TestApplicationList(t, &srv{}, clientFactory)
}

func TestApplicationData(t *testing.T) {
	integration.TestApplicationData(t, &srv{}, clientFactory)
}

func TestApplicationSubscribe(t *testing.T) {
	integration.TestApplicationSubscribe(t, &srv{}, clientFactory)
}

func TestApplicationUnsubscribe(t *testing.T) {
	integration.TestApplicationUnsubscribe(t, &srv{}, clientFactory)
}

func TestApplicationGet(t *testing.T) {
	integration.TestApplicationGet(t, &srv{}, clientFactory)
}
