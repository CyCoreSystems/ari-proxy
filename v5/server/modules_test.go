package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/v5/internal/integration"
)

func TestModulesData(t *testing.T) {
	integration.TestModulesData(t, &srv{})
}

func TestModulesLoad(t *testing.T) {
	integration.TestModulesLoad(t, &srv{})
}

func TestModulesReload(t *testing.T) {
	integration.TestModulesReload(t, &srv{})
}

func TestModulesUnload(t *testing.T) {
	integration.TestModulesUnload(t, &srv{})
}

func TestModulesList(t *testing.T) {
	integration.TestModulesList(t, &srv{})
}
