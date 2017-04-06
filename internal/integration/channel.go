package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestChannelData(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected ari.ChannelData
		expected.ID = "c1"
		expected.Name = "channe1"
		expected.State = "Up"

		m.Channel.On("Data", "c1").Return(&expected, nil)

		cd, err := cl.Channel().Data("c1")
		if err != nil {
			t.Errorf("Unexpected error in remote Data call %s", err)
		}

		if cd == nil || cd.ID != expected.ID || cd.Name != expected.Name || cd.State != expected.State {
			t.Errorf("Expected channel data %v, got %v", expected, cd)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Data", "c1")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = errors.New("Unknown error")

		m.Channel.On("Data", "c1").Return(nil, expected)

		cd, err := cl.Channel().Data("c1")
		if err == nil {
			t.Errorf("Expected error in remote Data call")
		}

		if cd != nil {
			t.Errorf("Expected channel data %v, got %v", nil, cd)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Data", "c1")
	})
}

func TestChannelAnswer(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Answer", "c1").Return(nil)

		err := cl.Channel().Answer("c1")
		if err != nil {
			t.Errorf("Unexpected error in remote Data call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Answer", "c1")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = errors.New("Unknown error")

		m.Channel.On("Answer", "c1").Return(expected)

		err := cl.Channel().Answer("c1")
		if err == nil {
			t.Errorf("Expected error in remote Answer call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Answer", "c1")
	})
}

func TestChannelBusy(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Busy", "c1").Return(nil)

		err := cl.Channel().Busy("c1")
		if err != nil {
			t.Errorf("Unexpected error in remote Busy call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Busy", "c1")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = errors.New("Unknown error")

		m.Channel.On("Busy", "c1").Return(expected)

		err := cl.Channel().Busy("c1")
		if err == nil {
			t.Errorf("Expected error in remote Busy call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Busy", "c1")
	})
}

func TestChannelCongestion(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Congestion", "c1").Return(nil)

		err := cl.Channel().Congestion("c1")
		if err != nil {
			t.Errorf("Unexpected error in remote Congestion call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Congestion", "c1")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = errors.New("Unknown error")

		m.Channel.On("Congestion", "c1").Return(expected)

		err := cl.Channel().Congestion("c1")
		if err == nil {
			t.Errorf("Expected error in remote Congestion call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Congestion", "c1")
	})
}

func TestChannelHangup(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("noReason", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Hangup", "c1", "").Return(nil)

		err := cl.Channel().Hangup("c1", "")
		if err != nil {
			t.Errorf("Unexpected error in remote Hangup call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", "c1", "")
	})

	runTest("aReason", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Hangup", "c1", "busy").Return(nil)

		err := cl.Channel().Hangup("c1", "busy")
		if err != nil {
			t.Errorf("Unexpected error in remote Hangup call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", "c1", "busy")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = errors.New("Unknown error")

		m.Channel.On("Hangup", "c1", "busy").Return(expected)

		err := cl.Channel().Hangup("c1", "busy")
		if err == nil {
			t.Errorf("Expected error in remote Hangup call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", "c1", "busy")
	})
}
