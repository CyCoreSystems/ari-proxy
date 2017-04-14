package server

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestMailboxList(t *testing.T) {
	integration.TestMailboxList(t, &srv{}, clientFactory)
}

func TestMailboxUpdate(t *testing.T) {
	integration.TestMailboxUpdate(t, &srv{}, clientFactory)
}

func TestMailboxDelete(t *testing.T) {
	integration.TestMailboxDelete(t, &srv{}, clientFactory)
}

func TestMailboxData(t *testing.T) {
	integration.TestMailboxData(t, &srv{}, clientFactory)
}
