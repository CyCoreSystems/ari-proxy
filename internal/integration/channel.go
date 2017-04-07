package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
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

func TestChannelList(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("empty", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("List").Return([]ari.ChannelHandle{}, nil)

		ret, err := cl.Channel().List()
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List")
	})

	runTest("nonEmpty", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var h1 = &mocks.ChannelHandle{}
		var h2 = &mocks.ChannelHandle{}

		h1.On("ID").Return("h1")
		h2.On("ID").Return("h2")

		m.Channel.On("List").Return([]ari.ChannelHandle{h1, h2}, nil)

		ret, err := cl.Channel().List()
		if err != nil {
			t.Errorf("Unexpected error in remote List call")
		}
		if len(ret) != 2 {
			t.Errorf("Expected return length to be 2, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List")
		h1.AssertCalled(t, "ID")
		h2.AssertCalled(t, "ID")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("List").Return(nil, errors.New("unknown error"))

		ret, err := cl.Channel().List()
		if err == nil {
			t.Errorf("Expected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List")
	})

}
