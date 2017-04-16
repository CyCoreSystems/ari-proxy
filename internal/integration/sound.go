package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

func TestSoundData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var sd ari.SoundData
		sd.ID = "s1"
		sd.Text = "text"

		m.Sound.On("Data", "s1").Return(&sd, nil)

		ret, err := cl.Sound().Data("s1")
		if err != nil {
			t.Errorf("Unexpected error in Sound Data: %s", err)
		}
		if ret == nil {
			t.Errorf("Expected data to be non-nil")
		} else {
			if ret.ID != sd.ID || ret.Text != sd.Text {
				t.Errorf("Expected '%v', got '%v'", sd, ret)
			}
		}

		m.Shutdown()

		m.Sound.AssertCalled(t, "Data", "s1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Sound.On("Data", "s1").Return(nil, errors.New("error"))

		_, err := cl.Sound().Data("s1")
		if err == nil {
			t.Errorf("Expected error in Sound Data: %s", err)
		}

		m.Shutdown()

		m.Sound.AssertCalled(t, "Data", "s1")
	})
}

func TestSoundList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		sh1 := &mocks.SoundHandle{}
		sh2 := &mocks.SoundHandle{}
		sh1.On("ID").Return("sh1")
		sh2.On("ID").Return("sh2")

		var filter map[string]string

		m.Sound.On("List", filter).Return([]ari.SoundHandle{sh1, sh2}, nil)

		ret, err := cl.Sound().List(nil)
		if err != nil {
			t.Errorf("Unexpected error in Sound List: %s", err)
		}
		if len(ret) != 2 {
			t.Errorf("Expected handle list length to be 2, got '%d'", len(ret))
		}

		m.Shutdown()

		m.Sound.AssertCalled(t, "List", filter)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var filter map[string]string

		m.Sound.On("List", filter).Return(nil, errors.New("error"))

		ret, err := cl.Sound().List(nil)
		if err == nil {
			t.Errorf("Expected error in Sound Data: %s", err)
		}
		if len(ret) != 0 {
			t.Errorf("Expected handle list length to be 0, got '%d'", len(ret))
		}
		m.Shutdown()

		m.Sound.AssertCalled(t, "List", filter)
	})
}
