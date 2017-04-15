package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

func TestDeviceData(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.DeviceStateData = "deviceData"
		m.DeviceState.On("Data", "d1").Return(&expected, nil)

		data, err := cl.DeviceState().Data("d1")
		if err != nil {
			t.Errorf("Error in remote device state data call: %s", err)
		}
		if data == nil || *data != expected {
			t.Errorf("Expected data '%s', got '%v'", expected, data)
		}

		m.DeviceState.AssertCalled(t, "Data", "d1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Data", "d1").Return(nil, errors.New("err"))

		data, err := cl.DeviceState().Data("d1")
		if err == nil {
			t.Errorf("Expected error in remote device state data call: %s", err)
		}
		if data != nil {
			t.Errorf("Expected data to be nil, got '%v'", *data)
		}

		m.DeviceState.AssertCalled(t, "Data", "d1")
	})
}

func TestDeviceDelete(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Delete", "d1").Return(nil)

		err := cl.DeviceState().Delete("d1")
		if err != nil {
			t.Errorf("Error in remote device state Delete call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Delete", "d1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Delete", "d1").Return(errors.New("err"))

		err := cl.DeviceState().Delete("d1")
		if err == nil {
			t.Errorf("Expected error in remote device state Delete call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Delete", "d1")
	})
}

func TestDeviceUpdate(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Update", "d1", "st1").Return(nil)

		err := cl.DeviceState().Update("d1", "st1")
		if err != nil {
			t.Errorf("Error in remote device state Update call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Update", "d1", "st1")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Update", "d1", "st1").Return(errors.New("err"))

		err := cl.DeviceState().Update("d1", "st1")
		if err == nil {
			t.Errorf("Expected error in remote device state Update call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Update", "d1", "st1")
	})
}

func TestDeviceList(t *testing.T, s Server, clientFactory ClientFactory) {
	runTest("ok", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		h1 := &mocks.DeviceStateHandle{}
		h2 := &mocks.DeviceStateHandle{}
		h1.On("ID").Return("h1")
		h2.On("ID").Return("h2")

		m.DeviceState.On("List").Return([]ari.DeviceStateHandle{h1, h2}, nil)

		list, err := cl.DeviceState().List()
		if err != nil {
			t.Errorf("Error in remote device state List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		h1.AssertCalled(t, "ID")
		h2.AssertCalled(t, "ID")
		m.DeviceState.AssertCalled(t, "List")
	})

	runTest("err", t, s, clientFactory, func(t *testing.T, m *mock, cl ari.Client) {

		m.DeviceState.On("List").Return([]ari.DeviceStateHandle{}, errors.New("error"))

		list, err := cl.DeviceState().List()
		if err == nil {
			t.Errorf("Expected error in remote device state List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.DeviceState.AssertCalled(t, "List")
	})
}
