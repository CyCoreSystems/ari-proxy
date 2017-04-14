package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestModulesData(t *testing.T) {
	integration.TestModulesData(t, &srv{}, clientFactory)
}

func TestModulesLoad(t *testing.T) {
	integration.TestModulesLoad(t, &srv{}, clientFactory)
}

func TestModulesReload(t *testing.T) {
	integration.TestModulesReload(t, &srv{}, clientFactory)
}

func TestModulesUnload(t *testing.T) {
	integration.TestModulesUnload(t, &srv{}, clientFactory)
}

func TestModulesList(t *testing.T) {
	integration.TestModulesList(t, &srv{}, clientFactory)
}
