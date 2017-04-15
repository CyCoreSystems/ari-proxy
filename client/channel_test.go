package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestChannelData(t *testing.T) {
	integration.TestChannelData(t, &srv{}, clientFactory)
}

func TestChannelAnswer(t *testing.T) {
	integration.TestChannelAnswer(t, &srv{}, clientFactory)
}

func TestChannelBusy(t *testing.T) {
	integration.TestChannelBusy(t, &srv{}, clientFactory)
}

func TestChannelCongestion(t *testing.T) {
	integration.TestChannelCongestion(t, &srv{}, clientFactory)
}

func TestChannelHangup(t *testing.T) {
	integration.TestChannelHangup(t, &srv{}, clientFactory)
}

func TestChannelList(t *testing.T) {
	integration.TestChannelList(t, &srv{}, clientFactory)
}

func TestChannelMute(t *testing.T) {
	integration.TestChannelMute(t, &srv{}, clientFactory)
}

func TestChannelUnmute(t *testing.T) {
	integration.TestChannelUnmute(t, &srv{}, clientFactory)
}

func TestChannelMOH(t *testing.T) {
	integration.TestChannelMOH(t, &srv{}, clientFactory)
}

func TestChannelStopMOH(t *testing.T) {
	integration.TestChannelStopMOH(t, &srv{}, clientFactory)
}

func TestChannelContinue(t *testing.T) {
	integration.TestChannelContinue(t, &srv{}, clientFactory)
}
