package server

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

func TestChannelCreate(t *testing.T) {
	integration.TestChannelCreate(t, &srv{}, clientFactory)
}

func TestChannelContinue(t *testing.T) {
	integration.TestChannelContinue(t, &srv{}, clientFactory)
}

func TestChannelDial(t *testing.T) {
	integration.TestChannelDial(t, &srv{}, clientFactory)
}

func TestChannelHold(t *testing.T) {
	integration.TestChannelHold(t, &srv{}, clientFactory)
}

func TestChannelStopHold(t *testing.T) {
	integration.TestChannelStopHold(t, &srv{}, clientFactory)
}

func TestChannelRing(t *testing.T) {
	integration.TestChannelRing(t, &srv{}, clientFactory)
}

func TestChannelStopRing(t *testing.T) {
	integration.TestChannelStopRing(t, &srv{}, clientFactory)
}

func TestChannelSilence(t *testing.T) {
	integration.TestChannelSilence(t, &srv{}, clientFactory)
}

func TestChannelStopSilence(t *testing.T) {
	integration.TestChannelStopSilence(t, &srv{}, clientFactory)
}

func TestChannelOriginate(t *testing.T) {
	integration.TestChannelOriginate(t, &srv{}, clientFactory)
}

func TestChannelPlay(t *testing.T) {
	integration.TestChannelPlay(t, &srv{}, clientFactory)
}

func TestChannelRecord(t *testing.T) {
	integration.TestChannelRecord(t, &srv{}, clientFactory)
}

func TestChannelSnoop(t *testing.T) {
	integration.TestChannelSnoop(t, &srv{}, clientFactory)
}

func TestChannelSendDTMF(t *testing.T) {
	integration.TestChannelSendDTMF(t, &srv{}, clientFactory)
}

func TestChannelVariableGet(t *testing.T) {
	integration.TestChannelVariableGet(t, &srv{}, clientFactory)
}

func TestChannelVariableSet(t *testing.T) {
	integration.TestChannelVariableSet(t, &srv{}, clientFactory)
}
