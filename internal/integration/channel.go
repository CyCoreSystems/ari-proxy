package integration

import (
	"errors"
	"testing"
	"time"

	rid "github.com/CyCoreSystems/ari-rid"
	"github.com/CyCoreSystems/ari/v5"
	tmock "github.com/stretchr/testify/mock"
)

var _ = tmock.Anything

var nonEmpty = tmock.MatchedBy(func(s string) bool { return s != "" })

func TestChannelData(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")
	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.ChannelData
		expected.ID = "c1"
		expected.Name = "channe1"
		expected.State = "Up"

		m.Channel.On("Data", key).Return(&expected, nil)

		cd, err := cl.Channel().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Data call %s", err)
		}

		if cd == nil || cd.ID != expected.ID || cd.Name != expected.Name || cd.State != expected.State {
			t.Errorf("Expected channel data %v, got %v", expected, cd)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Data", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("Unknown error")

		m.Channel.On("Data", key).Return(nil, expected)

		cd, err := cl.Channel().Data(key)
		if err == nil {
			t.Errorf("Expected error in remote Data call")
		}

		if cd != nil {
			t.Errorf("Expected channel data %v, got %v", nil, cd)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Data", key)
	})
}

func TestChannelAnswer(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Answer", key).Return(nil)

		err := cl.Channel().Answer(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Data call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Answer", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("Unknown error")

		m.Channel.On("Answer", key).Return(expected)

		err := cl.Channel().Answer(key)
		if err == nil {
			t.Errorf("Expected error in remote Answer call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Answer", key)
	})
}

func TestChannelBusy(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Busy", key).Return(nil)

		err := cl.Channel().Busy(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Busy call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Busy", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("Unknown error")

		m.Channel.On("Busy", key).Return(expected)

		err := cl.Channel().Busy(key)
		if err == nil {
			t.Errorf("Expected error in remote Busy call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Busy", key)
	})
}

func TestChannelCongestion(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Congestion", key).Return(nil)

		err := cl.Channel().Congestion(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Congestion call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Congestion", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("Unknown error")

		m.Channel.On("Congestion", key).Return(expected)

		err := cl.Channel().Congestion(key)
		if err == nil {
			t.Errorf("Expected error in remote Congestion call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Congestion", key)
	})
}

func TestChannelHangup(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("noReason", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Hangup", key, "").Return(nil)

		err := cl.Channel().Hangup(key, "")
		if err != nil {
			t.Errorf("Unexpected error in remote Hangup call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", key, "")
	})

	runTest("aReason", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Hangup", key, "busy").Return(nil)

		err := cl.Channel().Hangup(key, "busy")
		if err != nil {
			t.Errorf("Unexpected error in remote Hangup call %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", key, "busy")
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("Unknown error")

		m.Channel.On("Hangup", key, "busy").Return(expected)

		err := cl.Channel().Hangup(key, "busy")
		if err == nil {
			t.Errorf("Expected error in remote Hangup call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hangup", key, "busy")
	})
}

func TestChannelList(t *testing.T, s Server) {
	runTest("empty", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, nil)

		ret, err := cl.Channel().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote List call: %s", err)
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("nonEmpty", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		h1 := ari.NewKey(ari.ChannelKey, "h1")
		h2 := ari.NewKey(ari.ChannelKey, "h2")

		m.Channel.On("List", (*ari.Key)(nil)).Return([]*ari.Key{h1, h2}, nil)

		ret, err := cl.Channel().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in remote List call: %s", err)
		}
		if len(ret) != 2 {
			t.Errorf("Expected return length to be 2, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List", (*ari.Key)(nil))
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("List", (*ari.Key)(nil)).Return(nil, errors.New("unknown error"))

		ret, err := cl.Channel().List(nil)
		if err == nil {
			t.Errorf("Expected error in remote List call")
		}
		if len(ret) != 0 {
			t.Errorf("Expected return length to be 0, got %d", len(ret))
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "List", (*ari.Key)(nil))
	})
}

func TestChannelMute(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("both-ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", key, ari.DirectionBoth).Return(nil)

		err := cl.Channel().Mute(key, ari.DirectionBoth)
		if err != nil {
			t.Errorf("Unexpected error in remote Mute call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", key, ari.DirectionBoth)
	})

	runTest("both-err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", key, ari.DirectionBoth).Return(errors.New("error"))

		err := cl.Channel().Mute(key, ari.DirectionBoth)
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", key, ari.DirectionBoth)
	})

	runTest("none-ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", key, ari.DirectionNone).Return(nil)

		err := cl.Channel().Mute(key, ari.DirectionNone)
		if err != nil {
			t.Errorf("Unexpected error in remote Mute call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", key, ari.DirectionNone)
	})

	runTest("none-err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Mute", key, ari.DirectionNone).Return(errors.New("error"))

		err := cl.Channel().Mute(key, ari.DirectionNone)
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Mute", key, ari.DirectionNone)
	})
}

func TestChannelUnmute(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("both-ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", key, ari.DirectionBoth).Return(nil)

		err := cl.Channel().Unmute(key, ari.DirectionBoth)
		if err != nil {
			t.Errorf("Unexpected error in remote Unmute call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", key, ari.DirectionBoth)
	})

	runTest("both-err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", key, ari.DirectionBoth).Return(errors.New("error"))

		err := cl.Channel().Unmute(key, ari.DirectionBoth)
		if err == nil {
			t.Errorf("Expected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", key, ari.DirectionBoth)
	})

	runTest("none-ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", key, ari.DirectionNone).Return(nil)

		err := cl.Channel().Unmute(key, ari.DirectionNone)
		if err != nil {
			t.Errorf("Unexpected error in remote Unmute call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", key, ari.DirectionNone)
	})

	runTest("none-err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Unmute", key, ari.DirectionNone).Return(errors.New("error"))

		err := cl.Channel().Unmute(key, ari.DirectionNone)
		if err == nil {
			t.Errorf("Expected error in remote Unmute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Unmute", key, ari.DirectionNone)
	})
}

func TestChannelMOH(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("MOH", key, "music").Return(nil)

		err := cl.Channel().MOH(key, "music")
		if err != nil {
			t.Errorf("Unexpected error in remote MOH call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "MOH", key, "music")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("MOH", key, "music").Return(errors.New("error"))

		err := cl.Channel().MOH(key, "music")
		if err == nil {
			t.Errorf("Expected error in remote Mute call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "MOH", key, "music")
	})
}

func TestChannelStopMOH(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopMOH", key).Return(nil)

		err := cl.Channel().StopMOH(key)
		if err != nil {
			t.Errorf("Unexpected error in remote StopMOH call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopMOH", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopMOH", key).Return(errors.New("error"))

		err := cl.Channel().StopMOH(key)
		if err == nil {
			t.Errorf("Expected error in remote StopMOH call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopMOH", key)
	})
}

func TestChannelCreate(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		req := ari.ChannelCreateRequest{}
		req.App = "App"
		req.ChannelID = "1234"

		chkey := ari.NewKey(ari.ChannelKey, "ch2")
		expected := ari.NewChannelHandle(chkey, m.Channel, nil)

		m.Channel.On("Create", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req).Return(expected, nil)

		h, err := cl.Channel().Create(nil, req)
		if err != nil {
			t.Errorf("Unexpected error in remote Create call: %s", err)
		}
		if h == nil {
			t.Errorf("Expected non-nil channel handle")
		} else if h.ID() != req.ChannelID {
			t.Errorf("Expected handle id '%s', got '%s'", h.ID(), req.ChannelID)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Create", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, ari.ChannelCreateRequest{App: "App", ChannelID: "1234"})
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		req := ari.ChannelCreateRequest{}
		req.App = "App"
		req.ChannelID = "1234"

		m.Channel.On("Create", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req).Return(nil, errors.New("error"))

		h, err := cl.Channel().Create(nil, req)
		if err == nil {
			t.Errorf("Expected error in remote Create call")
		}
		if h != nil {
			t.Errorf("Expected nil channel handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Create", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, ari.ChannelCreateRequest{App: "App", ChannelID: "1234"})
	})
}

func TestChannelContinue(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Continue", key, "ctx1", "ext1", 0).Return(nil)

		err := cl.Channel().Continue(key, "ctx1", "ext1", 0)
		if err != nil {
			t.Errorf("Unexpected error in remote Continue call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Continue", key, "ctx1", "ext1", 0)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Continue", key, "ctx1", "ext1", 0).Return(errors.New("error"))

		err := cl.Channel().Continue(key, "ctx1", "ext1", 0)
		if err == nil {
			t.Errorf("Expected error in remote Continue call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Continue", key, "ctx1", "ext1", 0)
	})
}

func TestChannelDial(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Dial", key, "caller", 5*time.Second).Return(nil)

		err := cl.Channel().Dial(key, "caller", 5*time.Second)
		if err != nil {
			t.Errorf("Unexpected error in remote Dial call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Dial", key, "caller", 5*time.Second)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Dial", key, "caller", 5*time.Second).Return(errors.New("error"))

		err := cl.Channel().Dial(key, "caller", 5*time.Second)
		if err == nil {
			t.Errorf("Expected error in remote Dial call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Dial", key, "caller", 5*time.Second)
	})
}

func TestChannelHold(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Hold", key).Return(nil)

		err := cl.Channel().Hold(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Hold call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hold", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Hold", key).Return(errors.New("error"))

		err := cl.Channel().Hold(key)
		if err == nil {
			t.Errorf("Expected error in remote Hold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Hold", key)
	})
}

func TestChannelStopHold(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopHold", key).Return(nil)

		err := cl.Channel().StopHold(key)
		if err != nil {
			t.Errorf("Unexpected error in remote StopHold call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopHold", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopHold", key).Return(errors.New("error"))

		err := cl.Channel().StopHold(key)
		if err == nil {
			t.Errorf("Expected error in remote StopHold call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopHold", key)
	})
}

func TestChannelRing(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Ring", key).Return(nil)

		err := cl.Channel().Ring(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Ring call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Ring", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Ring", key).Return(errors.New("error"))

		err := cl.Channel().Ring(key)
		if err == nil {
			t.Errorf("Expected error in remote Ring call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Ring", key)
	})
}

func TestChannelSilence(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Silence", key).Return(nil)

		err := cl.Channel().Silence(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Silence call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Silence", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Silence", key).Return(errors.New("error"))

		err := cl.Channel().Silence(key)
		if err == nil {
			t.Errorf("Expected error in remote Silence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Silence", key)
	})
}

func TestChannelStopSilence(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopSilence", key).Return(nil)

		err := cl.Channel().StopSilence(key)
		if err != nil {
			t.Errorf("Unexpected error in remote StopSilence call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopSilence", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopSilence", key).Return(errors.New("error"))

		err := cl.Channel().StopSilence(key)
		if err == nil {
			t.Errorf("Expected error in remote StopSilence call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopSilence", key)
	})
}

func TestChannelStopRing(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopRing", key).Return(nil)

		err := cl.Channel().StopRing(key)
		if err != nil {
			t.Errorf("Unexpected error in remote StopRing call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopRing", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("StopRing", key).Return(errors.New("error"))

		err := cl.Channel().StopRing(key)
		if err == nil {
			t.Errorf("Expected error in remote StopRing call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "StopRing", key)
	})
}

func TestChannelOriginate(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := ari.NewChannelHandle(key, m.Channel, nil)

		var req ari.OriginateRequest
		req.App = "App"
		req.ChannelID = "1234"

		m.Channel.On("Originate", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req).Return(expected, nil)

		h, err := cl.Channel().Originate(nil, req)
		if err != nil {
			t.Errorf("Unexpected error in remote Originate call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != req.ChannelID {
			t.Errorf("Expected handle id '%s', got '%s'", h.ID(), req.ChannelID)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Originate", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var req ari.OriginateRequest
		req.App = "App"
		req.ChannelID = "1234"

		m.Channel.On("Originate", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req).Return(nil, errors.New("error"))

		h, err := cl.Channel().Originate(nil, req)
		if err == nil {
			t.Errorf("Expected error in remote Originate call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Originate", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, req)
	})
}

func TestChannelPlay(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")
	var cd ari.ChannelData
	cd.ID = "c1"
	cd.Name = "channe1"
	cd.State = "Up"

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		pbkey := ari.NewKey(ari.PlaybackKey, "pb1")
		expected := ari.NewPlaybackHandle(pbkey, m.Playback, nil)

		m.Channel.On("Data", key).Return(&cd, nil)
		m.Channel.On("Play", key, "playbackID", "sound:hello").Return(expected, nil)

		h, err := cl.Channel().Play(key, "playbackID", "sound:hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "playbackID" {
			t.Errorf("Expected handle id '%s', got '%s'", h.ID(), "playbackID")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Play", key, "playbackID", "sound:hello")
	})

	runTest("no-id", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		pbkey := ari.NewKey(ari.PlaybackKey, "pb1")
		expected := ari.NewPlaybackHandle(pbkey, m.Playback, nil)

		m.Channel.On("Data", key).Return(&cd, nil)
		m.Channel.On("Play", key, nonEmpty, "sound:hello").Return(expected, nil)

		h, err := cl.Channel().Play(key, "", "sound:hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Play call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Play", key, nonEmpty, "sound:hello")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("Data", key).Return(&cd, nil)
		m.Channel.On("Play", key, "playbackID", "sound:hello").Return(nil, errors.New("error"))

		h, err := cl.Channel().Play(key, "playbackID", "sound:hello")
		if err == nil {
			t.Errorf("Expected error in remote Play call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Play", key, "playbackID", "sound:hello")
	})
}

func TestChannelRecord(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.RecordingOptions

		lrkey := ari.NewKey(ari.LiveRecordingKey, "lrh1")
		expected := ari.NewLiveRecordingHandle(lrkey, m.LiveRecording, nil)

		m.Channel.On("Record", key, "recordid", opts).Return(expected, nil)

		h, err := cl.Channel().Record(key, "recordid", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "recordid" {
			t.Errorf("Expected handle id '%s', got '%s'", "recordid", h.ID())
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Record", key, "recordid", opts)
	})

	runTest("no-id", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.RecordingOptions

		lrkey := ari.NewKey(ari.LiveRecordingKey, "lrh1")
		expected := ari.NewLiveRecordingHandle(lrkey, m.LiveRecording, nil)

		m.Channel.On("Record", key, nonEmpty, opts).Return(expected, nil)

		h, err := cl.Channel().Record(key, "", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Record call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Record", key, nonEmpty, opts)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.RecordingOptions

		m.Channel.On("Record", key, "recordid", opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().Record(key, "recordid", nil)
		if err == nil {
			t.Errorf("Expected error in remote Record call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Record", key, "recordid", opts)
	})
}

func TestChannelSnoop(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.SnoopOptions

		chkey := ari.NewKey(ari.ChannelKey, "ch2")
		expected := ari.NewChannelHandle(chkey, m.Channel, nil)

		m.Channel.On("Snoop", key, "snoopID", opts).Return(expected, nil)

		h, err := cl.Channel().Snoop(key, "snoopID", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Snoop call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		} else if h.ID() != "snoopID" {
			t.Errorf("Expected handle id '%s', got '%s'", "snoopID", h.ID())
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Snoop", key, "snoopID", opts)
	})

	runTest("no-id", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.SnoopOptions

		chkey := ari.NewKey(ari.ChannelKey, "ch2")
		expected := ari.NewChannelHandle(chkey, m.Channel, nil)

		m.Channel.On("Snoop", key, tmock.Anything, opts).Return(expected, nil)

		h, err := cl.Channel().Snoop(key, "", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote Snoop call: %s", err)
		}
		if h == nil {
			t.Error("Expected non-nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Snoop", key, tmock.Anything, opts)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var opts *ari.SnoopOptions

		m.Channel.On("Snoop", key, "ch2", opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().Snoop(key, "ch2", nil)
		if err == nil {
			t.Errorf("Expected error in remote Snoop call")
		}
		if h != nil {
			t.Error("Expected nil handle")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "Snoop", key, "ch2", opts)
	})
}

func TestChannelExternalMedia(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		key := ari.NewKey(ari.ChannelKey, rid.New(rid.Channel))

		opts := ari.ExternalMediaOptions{
			ChannelID:    key.ID,
			App:          cl.ApplicationName(),
			ExternalHost: "localhost:1234",
			Format:       "slin16",
		}

		m.Channel.On("ExternalMedia", nil, opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().ExternalMedia(nil, opts)
		if err != nil {
			t.Errorf("error calling ExternalMedia: %v", err)
		} else if h.ID() != key.ID {
			t.Errorf("expected handle id '%s', got '%s'", key.ID, h.ID())
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "ExternalMedia", nil, opts)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		key := ari.NewKey(ari.ChannelKey, rid.New(rid.Channel))

		opts := ari.ExternalMediaOptions{
			ChannelID:    key.ID,
			App:          cl.ApplicationName(),
			ExternalHost: "localhost:1234",
			// Format: "slin16", // Format is required
		}

		m.Channel.On("ExternalMedia", nil, opts).Return(nil, errors.New("error"))

		h, err := cl.Channel().ExternalMedia(nil, opts)
		if err == nil {
			t.Error("expected error in ExternalMedia call, but got nil")
		}
		if h != nil {
			t.Errorf("expected nil channel handle, but got key with ID: %s", h.ID())
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "ExternalMedia", nil, opts)
	})
}

func TestChannelSendDTMF(t *testing.T, s Server) {
	key := ari.NewKey(ari.ChannelKey, "c1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("SendDTMF", key, "1", &ari.DTMFOptions{}).Return(nil)

		err := cl.Channel().SendDTMF(key, "1", nil)
		if err != nil {
			t.Errorf("Unexpected error in remote SendDTMF call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", key, "1", &ari.DTMFOptions{})
	})

	runTest("with-options", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		opts := &ari.DTMFOptions{
			After: 4 * time.Second,
		}

		m.Channel.On("SendDTMF", key, "1", opts).Return(nil)

		err := cl.Channel().SendDTMF(key, "1", opts)
		if err != nil {
			t.Errorf("Unexpected error in remote SendDTMF call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", key, "1", opts)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("SendDTMF", key, "1", &ari.DTMFOptions{}).Return(errors.New("err"))

		err := cl.Channel().SendDTMF(key, "1", nil)
		if err == nil {
			t.Errorf("Expected error in remote SendDTMF call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SendDTMF", key, "1", &ari.DTMFOptions{})
	})
}

func TestChannelVariableGet(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("GetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1").Return("value", nil)

		val, err := cl.Channel().GetVariable(nil, "v1")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Get call: %s", err)
		}
		if val != "value" {
			t.Errorf("Expected channel variable to be 'value', got '%s'", val)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "GetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("GetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1").Return("", errors.New("1"))

		val, err := cl.Channel().GetVariable(nil, "v1")
		if err == nil {
			t.Errorf("Expected error in remote Variables Get call")
		}
		if val != "" {
			t.Errorf("Expected empty value, got '%s'", val)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "GetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1")
	})
}

func TestChannelVariableSet(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("SetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1", "value").Return(nil)

		err := cl.Channel().SetVariable(nil, "v1", "value")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Set call: %s", err)
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1", "value")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Channel.On("SetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1", "value").Return(errors.New("error"))

		err := cl.Channel().SetVariable(nil, "v1", "value")
		if err == nil {
			t.Errorf("Expected error in remote Variables Set call")
		}

		m.Shutdown()

		m.Channel.AssertCalled(t, "SetVariable", &ari.Key{Kind: "", ID: "", Node: "", Dialog: "", App: ""}, "v1", "value")
	})
}
