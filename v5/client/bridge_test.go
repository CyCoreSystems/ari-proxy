package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/v5/internal/integration"
)

func TestBridgeCreate(t *testing.T) {
	integration.TestBridgeCreate(t, &srv{})
}

func TestBridgeList(t *testing.T) {
	integration.TestBridgeList(t, &srv{})
}

func TestBridgeData(t *testing.T) {
	integration.TestBridgeData(t, &srv{})
}

func TestBridgeAddChannel(t *testing.T) {
	integration.TestBridgeAddChannel(t, &srv{})
}

func TestBridgeRemoveChannel(t *testing.T) {
	integration.TestBridgeRemoveChannel(t, &srv{})
}

func TestBridgeDelete(t *testing.T) {
	integration.TestBridgeDelete(t, &srv{})
}

func TestBridgePlay(t *testing.T) {
	integration.TestBridgePlay(t, &srv{})
}

func TestBridgeRecord(t *testing.T) {
	integration.TestBridgeRecord(t, &srv{})
}
