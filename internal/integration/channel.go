package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
	tmock "github.com/stretchr/testify/mock"
)

var _ = tmock.Anything

var nonEmpty = tmock.MatchedBy(func(s string) bool { return s != "" })

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

func TestChannelMute(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("both-ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", "ch1", ari.DirectionBoth).Return(nil)

		err := cl.Channel().Mute("ch1", ari.DirectionBoth)
		if err != nil {
			t.Errorf("Unexpected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", "ch1", ari.DirectionBoth)
	})

	runTest("both-err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", "ch1", ari.DirectionBoth).Return(errors.New("error"))

		err := cl.Channel().Mute("ch1", ari.DirectionBoth)
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", "ch1", ari.DirectionBoth)
	})

	runTest("none-ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", "ch1", ari.DirectionNone).Return(nil)

		err := cl.Channel().Mute("ch1", ari.DirectionNone)
		if err != nil {
			t.Errorf("Unexpected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", "ch1", ari.DirectionNone)
	})

	runTest("none-err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", "ch1", ari.DirectionNone).Return(errors.New("error"))

		err := cl.Channel().Mute("ch1", ari.DirectionNone)
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", "ch1", ari.DirectionNone)
	})
}

func TestChannelUnmute(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("both-ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", "ch1", ari.DirectionBoth).Return(nil)

		err := cl.Channel().Unmute("ch1", ari.DirectionBoth)
		if err != nil {
			t.Errorf("Unexpected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", "ch1", ari.DirectionBoth)
	})

	runTest("both-err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", "ch1", ari.DirectionBoth).Return(errors.New("error"))

		err := cl.Channel().Unmute("ch1", ari.DirectionBoth)
		if err == nil {
			t.Errorf("Expected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", "ch1", ari.DirectionBoth)
	})

	runTest("none-ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", "ch1", ari.DirectionNone).Return(nil)

		err := cl.Channel().Unmute("ch1", ari.DirectionNone)
		if err != nil {
			t.Errorf("Unexpected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", "ch1", ari.DirectionNone)
	})

	runTest("none-err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", "ch1", ari.DirectionNone).Return(errors.New("error"))

		err := cl.Channel().Unmute("ch1", ari.DirectionNone)
		if err == nil {
			t.Errorf("Expected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", "ch1", ari.DirectionNone)
	})
}

func TestChannelMOH(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("MOH", "ch1", "music").Return(nil)

		err := cl.Channel().MOH("ch1", "music")
		if err != nil {
			t.Errorf("Unexpected error in remote MOH call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "MOH", "ch1", "music")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("MOH", "ch1", "music").Return(errors.New("error"))

		err := cl.Channel().MOH("ch1", "music")
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "MOH", "ch1", "music")
	})

}

func TestChannelStopMOH(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopMOH", "ch1").Return(nil)

		err := cl.Channel().StopMOH("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote StopMOH call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopMOH", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopMOH", "ch1").Return(errors.New("error"))

		err := cl.Channel().StopMOH("ch1")
		if err == nil {
			t.Errorf("Expected error in remote StopMOH call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopMOH", "ch1")
	})

}

func TestChannelCreate(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		req := ari.ChannelCreateRequest{}
		req.App = "App"

		expectedHandle := &mocks.ChannelHandle{}
		expectedHandle.On("ID").Return("ch1")

		m.Channel.On("Create", req).Return(expectedHandle, nil)

		h, err := cl.Channel().Create(req)
		if err != nil {
			t.Errorf("Unexpected error in remote Create call")
		}
		if h == nil {
			t.Errorf("Expected non-nil channel handle")
		} else if h.ID() != "ch1" {
			t.Errorf("Expected handle identifier 'ch1', got %s", h.ID())
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Create", ari.ChannelCreateRequest{App: "App"})
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		req := ari.ChannelCreateRequest{}
		req.App = "App"

		expectedHandle := &mocks.ChannelHandle{}
		expectedHandle.On("ID").Return("ch1")

		m.Channel.On("Create", req).Return(nil, errors.New("error"))

		h, err := cl.Channel().Create(req)
		if err == nil {
			t.Errorf("Expected error in remote Create call")
		}
		if h != nil {
			t.Errorf("Expected nil channel handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Create", ari.ChannelCreateRequest{App: "App"})
	})
}

func TestChannelContinue(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Continue", "ch1", "ctx1", "ext1", 0).Return(nil)

		err := cl.Channel().Continue("ch1", "ctx1", "ext1", 0)
		if err != nil {
			t.Errorf("Unexpected error in remote Continue call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Continue", "ch1", "ctx1", "ext1", 0)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Continue", "ch1", "ctx1", "ext1", 0).Return(errors.New("error"))

		err := cl.Channel().Continue("ch1", "ctx1", "ext1", 0)
		if err == nil {
			t.Errorf("Expected error in remote Continue call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Continue", "ch1", "ctx1", "ext1", 0)
	})
}

func TestChannelDial(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Dial", "ch1", "caller", 5*time.Second).Return(nil)

		err := cl.Channel().Dial("ch1", "caller", 5*time.Second)
		if err != nil {
			t.Errorf("Unexpected error in remote Dial call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Dial", "ch1", "caller", 5*time.Second)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Dial", "ch1", "caller", 5*time.Second).Return(errors.New("error"))

		err := cl.Channel().Dial("ch1", "caller", 5*time.Second)
		if err == nil {
			t.Errorf("Expected error in remote Dial call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Dial", "ch1", "caller", 5*time.Second)
	})
}

func TestChannelHold(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Hold", "ch1").Return(nil)

		err := cl.Channel().Hold("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote Hold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hold", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Hold", "ch1").Return(errors.New("error"))

		err := cl.Channel().Hold("ch1")
		if err == nil {
			t.Errorf("Expected error in remote Hold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hold", "ch1")
	})
}

func TestChannelStopHold(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopHold", "ch1").Return(nil)

		err := cl.Channel().StopHold("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote StopHold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopHold", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopHold", "ch1").Return(errors.New("error"))

		err := cl.Channel().StopHold("ch1")
		if err == nil {
			t.Errorf("Expected error in remote StopHold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopHold", "ch1")
	})
}

func TestChannelRing(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Ring", "ch1").Return(nil)

		err := cl.Channel().Ring("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote Ring call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Ring", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Ring", "ch1").Return(errors.New("error"))

		err := cl.Channel().Ring("ch1")
		if err == nil {
			t.Errorf("Expected error in remote Ring call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Ring", "ch1")
	})
}

func TestChannelSilence(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Silence", "ch1").Return(nil)

		err := cl.Channel().Silence("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote Silence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Silence", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Silence", "ch1").Return(errors.New("error"))

		err := cl.Channel().Silence("ch1")
		if err == nil {
			t.Errorf("Expected error in remote Silence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Silence", "ch1")
	})
}

func TestChannelStopSilence(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopSilence", "ch1").Return(nil)

		err := cl.Channel().StopSilence("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote StopSilence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopSilence", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopSilence", "ch1").Return(errors.New("error"))

		err := cl.Channel().StopSilence("ch1")
		if err == nil {
			t.Errorf("Expected error in remote StopSilence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopSilence", "ch1")
	})
}

func TestChannelStopRing(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopRing", "ch1").Return(nil)

		err := cl.Channel().StopRing("ch1")
		if err != nil {
			t.Errorf("Unexpected error in remote StopRing call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopRing", "ch1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("StopRing", "ch1").Return(errors.New("error"))

		err := cl.Channel().StopRing("ch1")
		if err == nil {
			t.Errorf("Expected error in remote StopRing call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopRing", "ch1")
	})
}

func TestChannelOriginate(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = &mocks.ChannelHandle{}
		expected.On("ID").Return("ch1")

		var req ari.OriginateRequest
		req.App = "App"

		m.Channel.On("Originate", req).Return(expected, nil)

		h, err := cl.Channel().Originate(req)
		if err != nil {
			t.Errorf("Unexpected error in remote Originate call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "ch1" {
			t.Errorf("Expected handle id 'ch1', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Originate", req)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var req ari.OriginateRequest
		req.App = "App"

		m.Channel.On("Originate", req).Return(nil, errors.New("error"))

		h, err := cl.Channel().Originate(req)
		if err == nil {
			t.Errorf("Expected error in remote Originate call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Originate", req)
	})
}

func TestChannelPlay(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = &mocks.PlaybackHandle{}
		expected.On("ID").Return("pb1")

		m.Channel.On("Play", "ch1", "playbackID", "sound:hello").Return(expected, nil)

		h, err := cl.Channel().Play("ch1", "playbackID", "sound:hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "pb1" {
			t.Errorf("Expected handle id 'pb1', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Play", "ch1", "playbackID", "sound:hello")
	})

	runTest("no-id", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var expected = &mocks.PlaybackHandle{}
		expected.On("ID").Return("pb1")

		m.Channel.On("Play", "ch1", nonEmpty, "sound:hello").Return(expected, nil)

		h, err := cl.Channel().Play("ch1", "", "sound:hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "pb1" {
			t.Errorf("Expected handle id 'pb1', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Play", "ch1", nonEmpty, "sound:hello")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("Play", "ch1", "playbackID", "sound:hello").Return(nil, errors.New("error"))

		h, err := cl.Channel().Play("ch1", "playbackID", "sound:hello")
		if err == nil {
			t.Errorf("Expected error in remote Play call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Play", "ch1", "playbackID", "sound:hello")
	})
}

func TestChannelRecord(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.RecordingOptions

		var expected = &mocks.LiveRecordingHandle{}
		expected.On("ID").Return("lrh1")

		m.Channel.On("Record", "ch1", "recordid", opts).Return(expected, nil)

		h, err := cl.Channel().Record("ch1", "recordid", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "lrh1" {
			t.Errorf("Expected handle id 'lrh1', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Record", "ch1", "recordid", opts)
	})

	runTest("no-id", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.RecordingOptions

		var expected = &mocks.LiveRecordingHandle{}
		expected.On("ID").Return("lrh1")

		m.Channel.On("Record", "ch1", nonEmpty, opts).Return(expected, nil)

		h, err := cl.Channel().Record("ch1", "", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "lrh1" {
			t.Errorf("Expected handle id 'lrh1', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Record", "ch1", nonEmpty, opts)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.RecordingOptions

		m.Channel.On("Record", "ch1", "recordid", opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().Record("ch1", "recordid", nil)
		if err == nil {
			t.Errorf("Expected error in remote Record call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Record", "ch1", "recordid", opts)
	})
}

func TestChannelSnoop(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.SnoopOptions

		var expected = &mocks.ChannelHandle{}
		expected.On("ID").Return("ch2")

		m.Channel.On("Snoop", "ch1", "snoopID", opts).Return(expected, nil)

		h, err := cl.Channel().Snoop("ch1", "snoopID", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Snoop call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "ch2" {
			t.Errorf("Expected handle id 'ch2', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Snoop", "ch1", "snoopID", opts)
	})

	runTest("no-id", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.SnoopOptions

		var expected = &mocks.ChannelHandle{}
		expected.On("ID").Return("ch2")

		m.Channel.On("Snoop", "ch1", "", opts).Return(expected, nil)

		h, err := cl.Channel().Snoop("ch1", "", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Snoop call")
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "ch2" {
			t.Errorf("Expected handle id 'ch2', got '%s'", h.ID())
		}

		m.Shutdown()

		expected.AssertCalled(t, "ID")
		m.Channel.AssertCalled(t, "Snoop", "ch1", "", opts)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		var opts *ari.SnoopOptions

		m.Channel.On("Snoop", "ch1", "ch2", opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().Snoop("ch1", "ch2", nil)
		if err == nil {
			t.Errorf("Expected error in remote Snoop call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Snoop", "ch1", "ch2", opts)
	})
}

func TestChannelSendDTMF(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("SendDTMF", "ch1", "1", &ari.DTMFOptions{}).Return(nil)

		err := cl.Channel().SendDTMF("ch1", "1", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote SendDTMF call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", "ch1", "1", &ari.DTMFOptions{})
	})

	runTest("with-options", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		opts := &ari.DTMFOptions{
			After: 4 * time.Second,
		}

		m.Channel.On("SendDTMF", "ch1", "1", opts).Return(nil)

		err := cl.Channel().SendDTMF("ch1", "1", opts)
		if err != nil {
			t.Errorf("Unexpected error in remote SendDTMF call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", "ch1", "1", opts)
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.Channel.On("SendDTMF", "ch1", "1", &ari.DTMFOptions{}).Return(errors.New("err"))

		err := cl.Channel().SendDTMF("ch1", "1", nil)
		if err == nil {
			t.Errorf("Expected error in remote SendDTMF call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", "ch1", "1", &ari.DTMFOptions{})
	})

}

func TestChannelVariableGet(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		mv := mocks.Variables{}
		mv.On("Get", "v1").Return("value", nil)
		m.Channel.On("Variables", "ch1").Return(&mv)

		val, err := cl.Channel().Variables("ch1").Get("v1")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Get call")
		}
		if val != "value" {
			t.Errorf("Expected channel variable to be 'value', got '%s'", val)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Variables", "ch1")
		mv.AssertCalled(t, "Get", "v1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		mv := mocks.Variables{}
		mv.On("Get", "v1").Return("", errors.New("error"))
		m.Channel.On("Variables", "ch1").Return(&mv)

		val, err := cl.Channel().Variables("ch1").Get("v1")
		if err == nil {
			t.Errorf("Expected error in remote Variables Get call")
		}
		if val != "" {
			t.Errorf("Expected empty value, got '%s'", val)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Variables", "ch1")
		mv.AssertCalled(t, "Get", "v1")
	})

}

func TestChannelVariableSet(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		mv := mocks.Variables{}
		mv.On("Set", "v1", "value").Return(nil)
		m.Channel.On("Variables", "ch1").Return(&mv)

		err := cl.Channel().Variables("ch1").Set("v1", "value")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Set call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Variables", "ch1")
		mv.AssertCalled(t, "Set", "v1", "value")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		mv := mocks.Variables{}
		mv.On("Set", "v1", "value").Return(errors.New("error"))
		m.Channel.On("Variables", "ch1").Return(&mv)

		err := cl.Channel().Variables("ch1").Set("v1", "value")
		if err == nil {
			t.Errorf("Expected error in remote Variables Set call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Variables", "ch1")
		mv.AssertCalled(t, "Set", "v1", "value")
	})

}
