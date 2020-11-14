package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/arimocks"
	"github.com/rotisserie/eris"
)

func TestAsteriskInfo(t *testing.T, s Server) {
	runTest("noFilter", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		m.Asterisk.On("Info", ari.NodeKey("asdf", "1")).Return(&ai, nil)

		ret, err := cl.Asterisk().Info(ari.NodeKey("asdf", "1"))
		if err != nil {
			t.Errorf("Unexpected error in remote Info call: %v", err)
		}
		if ret == nil || ret.SystemInfo.EntityID != ai.SystemInfo.EntityID {
			t.Errorf("Expected asterisk info %v, got %v", ai, ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", ari.NodeKey("asdf", "1"))
	})

	runTest("noFilterError", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := eris.New("unknown error")

		m.Asterisk.On("Info", ari.NodeKey("asdf", "1")).Return(nil, expected)

		ret, err := cl.Asterisk().Info(ari.NodeKey("asdf", "1"))
		if err == nil || eris.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != nil {
			t.Errorf("Expected nil ret, got '%v'", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", ari.NodeKey("asdf", "1"))
	})
}

func TestAsteriskVariablesGet(t *testing.T, s Server) {
	key := ari.NewKey(ari.VariableKey, "s")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := arimocks.AsteriskVariables{}
		mv.On("Get", key).Return("hello", nil)
		m.Asterisk.On("Variables").Return(&mv)

		ret, err := cl.Asterisk().Variables().Get(key)
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Get call: %v", err)
		}
		if ret != "hello" {
			t.Errorf("Expected Variables Get return %v, got %v", "hello", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Get", key)
	})

	runTest("error", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		expected := eris.New("unknown error")

		mv := arimocks.AsteriskVariables{}
		mv.On("Get", key).Return("", expected)

		m.Asterisk.On("Variables").Return(&mv)

		ret, err := cl.Asterisk().Variables().Get(key)
		if err == nil || eris.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != "" {
			t.Errorf("Expected Variables Get return %v, got %v", "", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Get", key)
	})
}

func TestAsteriskVariablesSet(t *testing.T, s Server) {
	key := ari.NewKey(ari.VariableKey, "s")

	runTest("simple", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := arimocks.AsteriskVariables{}
		m.Asterisk.On("Variables").Return(&mv)
		mv.On("Set", key, "hello").Return(nil)

		err := cl.Asterisk().Variables().Set(key, "hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Set call: %v", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Set", key, "hello")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := arimocks.AsteriskVariables{}
		m.Asterisk.On("Variables").Return(&mv)
		mv.On("Set", key, "hello").Return(eris.New("error"))

		err := cl.Asterisk().Variables().Set(key, "hello")
		if err == nil {
			t.Errorf("Expected error in remote Variables Set call: %v", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Set", key, "hello")
	})
}
