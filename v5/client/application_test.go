package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/v5/internal/integration"
)

func TestApplicationList(t *testing.T) {
	integration.TestApplicationList(t, &srv{})
}

func TestApplicationData(t *testing.T) {
	integration.TestApplicationData(t, &srv{})
}

func TestApplicationSubscribe(t *testing.T) {
	integration.TestApplicationSubscribe(t, &srv{})
}

func TestApplicationUnsubscribe(t *testing.T) {
	integration.TestApplicationUnsubscribe(t, &srv{})
}

func TestApplicationGet(t *testing.T) {
	integration.TestApplicationGet(t, &srv{})
}
