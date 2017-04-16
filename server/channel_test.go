package server

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

func TestChannelCreate(t *testing.T) {
	integration.TestChannelCreate(t, &srv{})
}

func TestChannelContinue(t *testing.T) {
	integration.TestChannelContinue(t, &srv{})
}

func TestChannelDial(t *testing.T) {
	integration.TestChannelDial(t, &srv{})
}

func TestChannelHold(t *testing.T) {
	integration.TestChannelHold(t, &srv{})
}

func TestChannelStopHold(t *testing.T) {
	integration.TestChannelStopHold(t, &srv{})
}

func TestChannelRing(t *testing.T) {
	integration.TestChannelRing(t, &srv{})
}

func TestChannelStopRing(t *testing.T) {
	integration.TestChannelStopRing(t, &srv{})
}

func TestChannelSilence(t *testing.T) {
	integration.TestChannelSilence(t, &srv{})
}

func TestChannelStopSilence(t *testing.T) {
	integration.TestChannelStopSilence(t, &srv{})
}

func TestChannelOriginate(t *testing.T) {
	integration.TestChannelOriginate(t, &srv{})
}

func TestChannelPlay(t *testing.T) {
	integration.TestChannelPlay(t, &srv{})
}

func TestChannelRecord(t *testing.T) {
	integration.TestChannelRecord(t, &srv{})
}

func TestChannelSnoop(t *testing.T) {
	integration.TestChannelSnoop(t, &srv{})
}

func TestChannelSendDTMF(t *testing.T) {
	integration.TestChannelSendDTMF(t, &srv{})
}

func TestChannelVariableGet(t *testing.T) {
	integration.TestChannelVariableGet(t, &srv{})
}

func TestChannelVariableSet(t *testing.T) {
	integration.TestChannelVariableSet(t, &srv{})
}
