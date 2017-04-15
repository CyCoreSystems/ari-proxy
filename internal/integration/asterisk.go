package integration

import (
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
	"github.com/pkg/errors"
)

func TestAsteriskInfo(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("noFilter", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		m.Asterisk.On("Info", "").Return(&ai, nil)

		ret, err := cl.Asterisk().Info("")
		if err != nil {
			t.Errorf("Unexpected error in remote Info call: %v", err)
		}
		if ret == nil || ret.SystemInfo.EntityID != ai.SystemInfo.EntityID {
			t.Errorf("Expected asterisk info %v, got %v", ai, ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", "")
	})

	runTest("noFilterError", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Asterisk.On("Info", "").Return(nil, expected)

		ret, err := cl.Asterisk().Info("")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != nil {
			t.Errorf("Expected nil ret, got '%v'", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", "")
	})

	runTest("withFilter", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		m.Asterisk.On("Info", "").Return(&ai, nil)

		ret, err := cl.Asterisk().Info("filter")
		if err != nil {
			t.Errorf("Unexpected error in remote Info call: %v", err)
		}
		if ret == nil || ret.SystemInfo.EntityID != ai.SystemInfo.EntityID {
			t.Errorf("Expected asterisk info %v, got %v", ai, ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", "") // filter gets stripped since it isn't supported
	})

	runTest("withFilterError", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		expected := errors.New("unknown error")

		m.Asterisk.On("Info", "").Return(nil, expected)

		ret, err := cl.Asterisk().Info("filter")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != nil {
			t.Errorf("Expected nil ret, got '%v'", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Info", "") // filter gets stripped since it isn't supported
	})
}

func TestAsteriskVariablesGet(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := mocks.Variables{}
		mv.On("Get", "s").Return("hello", nil)
		m.Asterisk.On("Variables").Return(&mv)

		ret, err := cl.Asterisk().Variables().Get("s")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Get call: %v", err)
		}
		if ret != "hello" {
			t.Errorf("Expected Variables Get return %v, got %v", "hello", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Get", "s")
	})

	runTest("error", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var expected = errors.New("unknown error")

		mv := mocks.Variables{}
		mv.On("Get", "s").Return("", expected)

		m.Asterisk.On("Variables").Return(&mv)

		ret, err := cl.Asterisk().Variables().Get("s")
		if err == nil || errors.Cause(err).Error() != expected.Error() {
			t.Errorf("Expected error '%v', got '%v'", expected, err)
		}
		if ret != "" {
			t.Errorf("Expected Variables Get return %v, got %v", "", ret)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Get", "s")
	})
}

func TestAsteriskVariablesSet(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("simple", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := mocks.Variables{}
		m.Asterisk.On("Variables").Return(&mv)
		mv.On("Set", "s", "hello").Return(nil)

		err := cl.Asterisk().Variables().Set("s", "hello")
		if err != nil {
			t.Errorf("Unexpected error in remote Variables Set call: %v", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Set", "s", "hello")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var ai ari.AsteriskInfo
		ai.SystemInfo.EntityID = "1"

		mv := mocks.Variables{}
		m.Asterisk.On("Variables").Return(&mv)
		mv.On("Set", "s", "hello").Return(errors.New("error"))

		err := cl.Asterisk().Variables().Set("s", "hello")
		if err == nil {
			t.Errorf("Expected error in remote Variables Set call: %v", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Variables")
		mv.AssertCalled(t, "Set", "s", "hello")
	})
}
