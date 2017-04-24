package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestDeviceData(t *testing.T, s Server) {
	key := ari.NewKey(ari.DeviceStateKey, "d1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		var expected ari.DeviceStateData
		expected.State = "deviceData1"
		expected.Key = ari.NewKey(ari.DeviceStateKey, "d1")

		m.DeviceState.On("Data", key).Return(&expected, nil)

		data, err := cl.DeviceState().Data(key)
		if err != nil {
			t.Errorf("Error in remote device state data call: %s", err)
		}
		if data == nil || data.State != expected.State {
			t.Errorf("Expected data '%s', got '%v'", expected, data)
		}

		m.DeviceState.AssertCalled(t, "Data", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Data", key).Return(nil, errors.New("err"))

		data, err := cl.DeviceState().Data(key)
		if err == nil {
			t.Errorf("Expected error in remote device state data call: %s", err)
		}
		if data != nil {
			t.Errorf("Expected data to be nil, got '%v'", *data)
		}

		m.DeviceState.AssertCalled(t, "Data", key)
	})
}

func TestDeviceDelete(t *testing.T, s Server) {
	key := ari.NewKey(ari.DeviceStateKey, "d1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Delete", key).Return(nil)

		err := cl.DeviceState().Delete(key)
		if err != nil {
			t.Errorf("Error in remote device state Delete call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Delete", key)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Delete", key).Return(errors.New("err"))

		err := cl.DeviceState().Delete(key)
		if err == nil {
			t.Errorf("Expected error in remote device state Delete call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Delete", key)
	})
}

func TestDeviceUpdate(t *testing.T, s Server) {
	key := ari.NewKey(ari.DeviceStateKey, "d1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Update", "d1", "st1").Return(nil)

		err := cl.DeviceState().Update(key, "st1")
		if err != nil {
			t.Errorf("Error in remote device state Update call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Update", key, "st1")
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.DeviceState.On("Update", key, "st1").Return(errors.New("err"))

		err := cl.DeviceState().Update(key, "st1")
		if err == nil {
			t.Errorf("Expected error in remote device state Update call: %s", err)
		}

		m.DeviceState.AssertCalled(t, "Update", key, "st1")
	})
}

func TestDeviceList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		h1 := ari.NewKey(ari.DeviceStateKey, "h1")
		h2 := ari.NewKey(ari.DeviceStateKey, "h2")

		m.DeviceState.On("List", nil).Return([]*ari.Key{h1, h2}, nil)

		list, err := cl.DeviceState().List(nil)
		if err != nil {
			t.Errorf("Error in remote device state List call: %s", err)
		}
		if len(list) != 2 {
			t.Errorf("Expected list of length 2, got %d", len(list))
		}

		m.DeviceState.AssertCalled(t, "List", nil)
	})

	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.DeviceState.On("List").Return([]*ari.Key{}, errors.New("error"))

		list, err := cl.DeviceState().List(nil)
		if err == nil {
			t.Errorf("Expected error in remote device state List call")
		}
		if len(list) != 0 {
			t.Errorf("Expected list of length 0, got %d", len(list))
		}

		m.DeviceState.AssertCalled(t, "List", nil)
	})
}
