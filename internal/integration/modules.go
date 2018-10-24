package integration

import (
	"errors"
	"testing"

	"github.com/CyCoreSystems/ari"
)

func TestModulesLoad(t *testing.T, s Server) {
	key := ari.NewKey(ari.ModuleKey, "m1")
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Load", key).Return(nil)

		if err := cl.Asterisk().Modules().Load(key); err != nil {
			t.Errorf("Unexpected error in module load: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Load", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Load", key).Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Load(key); err == nil {
			t.Errorf("Expected error in module load")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Load", key)
	})
}

func TestModulesUnload(t *testing.T, s Server) {
	key := ari.NewKey(ari.ModuleKey, "m1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Unload", key).Return(nil)

		if err := cl.Asterisk().Modules().Unload(key); err != nil {
			t.Errorf("Unexpected error in module Unload: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Unload", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Unload", key).Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Unload(key); err == nil {
			t.Errorf("Expected error in module Unload")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Unload", key)
	})
}

func TestModulesReload(t *testing.T, s Server) {
	key := ari.NewKey(ari.ModuleKey, "m1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Reload", key).Return(nil)

		if err := cl.Asterisk().Modules().Reload(key); err != nil {
			t.Errorf("Unexpected error in module Reload: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Reload", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Reload", key).Return(errors.New("error"))

		if err := cl.Asterisk().Modules().Reload(key); err == nil {
			t.Errorf("Expected error in module Reload")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Reload", key)
	})
}

func TestModulesData(t *testing.T, s Server) {
	key := ari.NewKey(ari.ModuleKey, "m1")

	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		var d ari.ModuleData
		d.Description = "Desc"
		d.Name = "name"

		m.Modules.On("Data", key).Return(&d, nil)

		ret, err := cl.Asterisk().Modules().Data(key)
		if err != nil {
			t.Errorf("Unexpected error in module Data: %s", err)
		}
		if ret == nil {
			t.Errorf("Expected module data to be non-nil")
		} else {
			if ret.Description != d.Description {
				t.Errorf("description mismatch: %v %v", ret.Description, d.Description)
			}
			if ret.Name != d.Name {
				t.Errorf("name mismatch: expected '%v', got '%v'", d, ret)
			}
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Data", key)
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {

		m.Modules.On("Data", key).Return(nil, errors.New("error"))

		_, err := cl.Asterisk().Modules().Data(key)
		if err == nil {
			t.Errorf("Expected error in module Data: %s", err)
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "Data", key)
	})
}

func TestModulesList(t *testing.T, s Server) {
	runTest("ok", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m1 := ari.NewKey(ari.ModuleKey, "m1")
		m2 := ari.NewKey(ari.ModuleKey, "m2")

		m.Modules.On("List", (*ari.Key)(nil)).Return([]*ari.Key{m1, m2}, nil)

		ret, err := cl.Asterisk().Modules().List(nil)
		if err != nil {
			t.Errorf("Unepected error in module List: %s", err)
		}
		if len(ret) != 2 {
			t.Errorf("Expected handle list length of size '2', got '%d'", len(ret))
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "List", (*ari.Key)(nil))
	})
	runTest("err", t, s, func(t *testing.T, m *mock, cl ari.Client) {
		m.Modules.On("List", (*ari.Key)(nil)).Return([]*ari.Key{}, errors.New("error"))

		_, err := cl.Asterisk().Modules().List(nil)
		if err == nil {
			t.Errorf("Expected error in module List")
		}

		m.Shutdown()

		m.Asterisk.AssertCalled(t, "Modules")
		m.Modules.AssertCalled(t, "List", (*ari.Key)(nil))
	})
}
