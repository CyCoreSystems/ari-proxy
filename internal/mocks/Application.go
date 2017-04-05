package mocks

import ari "github.com/CyCoreSystems/ari"
import mock "github.com/stretchr/testify/mock"

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// Data provides a mock function with given fields: name
func (_m *Application) Data(name string) (*ari.ApplicationData, error) {
	ret := _m.Called(name)

	var r0 *ari.ApplicationData
	if rf, ok := ret.Get(0).(func(string) *ari.ApplicationData); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.ApplicationData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: name
func (_m *Application) Get(name string) ari.ApplicationHandle {
	ret := _m.Called(name)

	var r0 ari.ApplicationHandle
	if rf, ok := ret.Get(0).(func(string) ari.ApplicationHandle); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ari.ApplicationHandle)
		}
	}

	return r0
}

// List provides a mock function with given fields:
func (_m *Application) List() ([]ari.ApplicationHandle, error) {
	ret := _m.Called()

	var r0 []ari.ApplicationHandle
	if rf, ok := ret.Get(0).(func() []ari.ApplicationHandle); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ari.ApplicationHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: name, eventSource
func (_m *Application) Subscribe(name string, eventSource string) error {
	ret := _m.Called(name, eventSource)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, eventSource)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with given fields: name, eventSource
func (_m *Application) Unsubscribe(name string, eventSource string) error {
	ret := _m.Called(name, eventSource)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, eventSource)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
