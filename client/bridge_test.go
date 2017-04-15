package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestBridgeCreate(t *testing.T) {
	integration.TestBridgeCreate(t, &srv{}, clientFactory)
}

func TestBridgeList(t *testing.T) {
	integration.TestBridgeList(t, &srv{}, clientFactory)
}

func TestBridgeData(t *testing.T) {
	integration.TestBridgeData(t, &srv{}, clientFactory)
}

func TestBridgeAddChannel(t *testing.T) {
	integration.TestBridgeAddChannel(t, &srv{}, clientFactory)
}

func TestBridgeRemoveChannel(t *testing.T) {
	integration.TestBridgeRemoveChannel(t, &srv{}, clientFactory)
}

func TestBridgeDelete(t *testing.T) {
	integration.TestBridgeDelete(t, &srv{}, clientFactory)
}

func TestBridgePlay(t *testing.T) {
	integration.TestBridgePlay(t, &srv{}, clientFactory)
}

func TestBridgeRecord(t *testing.T) {
	integration.TestBridgeRecord(t, &srv{}, clientFactory)
}
