package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestMailboxList(t *testing.T) {
	integration.TestMailboxList(t, &srv{})
}

func TestMailboxUpdate(t *testing.T) {
	integration.TestMailboxUpdate(t, &srv{})
}

func TestMailboxDelete(t *testing.T) {
	integration.TestMailboxDelete(t, &srv{})
}

func TestMailboxData(t *testing.T) {
	integration.TestMailboxData(t, &srv{})
}
