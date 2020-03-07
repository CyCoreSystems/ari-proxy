package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/arimocks"
	tmock "github.com/stretchr/testify/mock"
)

var _ = tmock.Anything

func TestConfigData(t *testing.T, s Server) {
	key := ari.NewKey("config", ari.ConfigID("c1", "o1", "id1"))

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.ConfigData
		expected.Class = "c1"
		expected.Key = key
		expected.Type = "o1"
		expected.Fields = []ari.ConfigTuple{
			ari.ConfigTuple{
				Value:     "v1",
				Attribute: "a1",
			},
		}

		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Data", key).Return(&expected, nil)

		data, err := cl.Asterisk().Config().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in remove config data call: %s", err)
		}
		if data == nil {
			t.Errorf("Expected non-nil data")
		} else {
			failed := false
			failed = failed || data.Class != expected.Class
			failed = failed || data.ID() != expected.ID()
			failed = failed || data.Type != expected.Type
			failed = failed || len(data.Fields) != len(expected.Fields)
			for idx := range data.Fields {
				failed = failed || data.Fields[idx].Attribute != expected.Fields[idx].Attribute
				failed = failed || data.Fields[idx].Value != expected.Fields[idx].Value
			}
			if failed {
				t.Errorf("Expected config data '%v', got '%v'", expected, data)
			}
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Data", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Data", key).Return(nil, errors.New("error"))

		data, err := cl.Asterisk().Config().Data(key)
		if err == nil {
			t.Errorf("Expected error in remove config data call")
		}
		if data != nil {
			t.Errorf("Expected nil data")
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Data", key)
	})
}

func TestConfigDelete(t *testing.T, s Server) {
	key := ari.NewKey("config", ari.ConfigID("c1", "o1", "id1"))

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Delete", key).Return(nil)

		err := cl.Asterisk().Config().Delete(key)
		if err != nil {
			t.Errorf("Unexpected error in remove config delete call: %s", err)
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Delete", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Delete", key).Return(errors.New("error"))

		err := cl.Asterisk().Config().Delete(key)
		if err == nil {
			t.Errorf("Expected error in remove config delete call")
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Delete", key)
	})
}

func TestConfigUpdate(t *testing.T, s Server) {
	key := ari.NewKey("config", ari.ConfigID("c1", "o1", "id1"))

	tuples := []ari.ConfigTuple{
		ari.ConfigTuple{
			Value:     "v1",
			Attribute: "a1",
		},
		ari.ConfigTuple{
			Value:     "v2",
			Attribute: "a2",
		},
	}

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Update", key, tuples).Return(nil)

		err := cl.Asterisk().Config().Update(key, tuples)
		if err != nil {
			t.Errorf("Unexpected error in remove config Update call: %s", err)
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Update", key, tuples)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		cfg := arimocks.Config{}
		m.Asterisk.On("Config").Return(&cfg)
		cfg.On("Update", key, tuples).Return(errors.New("error"))

		err := cl.Asterisk().Config().Update(key, tuples)
		if err == nil {
			t.Errorf("Expected error in remove config Update call")
		}

		m.Asterisk.AssertCalled(t, "Config")
		cfg.AssertCalled(t, "Update", key, tuples)
	})
}
