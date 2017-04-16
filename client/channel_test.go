package client

import (
	"testing"

	"github.com/CyCoreSystems/ari-proxy/internal/integration"
)

func TestChannelData(t *testing.T) {
	integration.TestChannelData(t, &srv{})
}

func TestChannelAnswer(t *testing.T) {
	integration.TestChannelAnswer(t, &srv{})
}

func TestChannelBusy(t *testing.T) {
	integration.TestChannelBusy(t, &srv{})
}

func TestChannelCongestion(t *testing.T) {
	integration.TestChannelCongestion(t, &srv{})
}

func TestChannelHangup(t *testing.T) {
	integration.TestChannelHangup(t, &srv{})
}

func TestChannelList(t *testing.T) {
	integration.TestChannelList(t, &srv{})
}

func TestChannelMute(t *testing.T) {
	integration.TestChannelMute(t, &srv{})
}

func TestChannelUnmute(t *testing.T) {
	integration.TestChannelUnmute(t, &srv{})
}

func TestChannelMOH(t *testing.T) {
	integration.TestChannelMOH(t, &srv{})
}

func TestChannelStopMOH(t *testing.T) {
	integration.TestChannelStopMOH(t, &srv{})
}

func TestChannelContinue(t *testing.T) {
	integration.TestChannelContinue(t, &srv{})
}
