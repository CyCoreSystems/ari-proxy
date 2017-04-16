package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

func TestModulesLoad(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Load", "m1").Return(nil)

		if err := cl.Asterisk().Modules().Load("m1"); err != nil {
			t.Errorf("Unexpected error in module load: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Load", "m1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Load", "m1").Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Load("m1"); err == nil {
			t.Errorf("Expected error in module load")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Load", "m1")
	})
}

func TestModulesUnload(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Unload", "m1").Return(nil)

		if err := cl.Asterisk().Modules().Unload("m1"); err != nil {
			t.Errorf("Unexpected error in module Unload: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Unload", "m1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Unload", "m1").Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Unload("m1"); err == nil {
			t.Errorf("Expected error in module Unload")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Unload", "m1")
	})
}

func TestModulesReload(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Reload", "m1").Return(nil)

		if err := cl.Asterisk().Modules().Reload("m1"); err != nil {
			t.Errorf("Unexpected error in module Reload: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Reload", "m1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Reload", "m1").Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Reload("m1"); err == nil {
			t.Errorf("Expected error in module Reload")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Reload", "m1")
	})
}

func TestModulesData(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var d ari.ModuleData
		d.Description = "Desc"
		d.Name = "name"

		m.Modules.On("Data", "m1").Return(&d, nil)

		ret, err := cl.Asterisk().Modules().Data("m1")
		if err != nil {
			t.Errorf("Unexpected error in module Data: %s", err)
		}
		if ret == nil {
			t.Errorf("Expected module data to be non-nil")
		} else {
			failed := ret.Description != d.Description
			failed = ret.Name != d.Name
			if failed {
				t.Errorf("Expected '%v', got '%v'", d, ret)
			}
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Data", "m1")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Data", "m1").Return(nil, errors.New("error"))

		_, err := cl.Asterisk().Modules().Data("m1")
		if err == nil {
			t.Errorf("Expected error in module Data: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Data", "m1")
	})
}

func TestModulesList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m1 := mocks.ModuleHandle{}
		m1.On("ID").Return("m1")
		m2 := mocks.ModuleHandle{}
		m2.On("ID").Return("m2")

		m.Modules.On("List").Return([]ari.ModuleHandle{&m1, &m2}, nil)

		ret, err := cl.Asterisk().Modules().List()
		if err != nil {
			t.Errorf("Unepected error in module List: %s", err)
		}
		if len(ret) != 2 {
			t.Errorf("Expected handle list length of size '2', got '%d'", len(ret))
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "List")
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Modules.On("List").Return([]ari.ModuleHandle{}, errors.New("error"))

		_, err := cl.Asterisk().Modules().List()
		if err == nil {
			t.Errorf("Expected error in module List")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "List")
	})
}
