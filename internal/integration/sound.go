package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestSoundData(t *testing.T, s Server) {
	key := ari.NewKey(ari.SoundKey, "s1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var sd ari.SoundData
		sd.ID = "s1"
		sd.Text = "text"

		m.Sound.On("Data", key).Return(&sd, nil)

		ret, err := cl.Sound().Data(key)
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

		m.Sound.AssertCalled(t, "Data", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Sound.On("Data", key).Return(nil, errors.New("error"))

		_, err := cl.Sound().Data(key)
		if err == nil {
			t.Errorf("Expected error in Sound Data: %s", err)
		}

		m.Shutdown()

		m.Sound.AssertCalled(t, "Data", key)
	})
}

func TestSoundList(t *testing.T, s Server) {

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		sh1 := ari.NewKey(ari.SoundKey, "sh1")
		sh2 := ari.NewKey(ari.SoundKey, "sh2")

		var filter map[string]string

		m.Sound.On("List", filter, nil).Return([]*ari.Key{sh1, sh2}, nil)

		ret, err := cl.Sound().List(nil, nil)
		if err != nil {
			t.Errorf("Unexpected error in Sound List: %s", err)
		}
		if len(ret) != 2 {
			t.Errorf("Expected handle list length to be 2, got '%d'", len(ret))
		}

		m.Shutdown()

		m.Sound.AssertCalled(t, "List", filter, nil)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var filter map[string]string

		m.Sound.On("List", filter, nil).Return(nil, errors.New("error"))

		ret, err := cl.Sound().List(nil, nil)
		if err == nil {
			t.Errorf("Expected error in Sound Data: %s", err)
		}
		if len(ret) != 0 {
			t.Errorf("Expected handle list length to be 0, got '%d'", len(ret))
		}
		m.Shutdown()

		m.Sound.AssertCalled(t, "List", filter, nil)
	})
}
