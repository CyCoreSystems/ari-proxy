package mocks

import ari "github.com/CyCoreSystems/ari"
import mock "github.com/stretchr/testify/mock"

// PlaybackHandle is an autogenerated mock type for the PlaybackHandle type
type PlaybackHandle struct {
	mock.Mock
}

// Control provides a mock function with given fields: op
func (_m *PlaybackHandle) Control(op string) error {
	ret := _m.Called(op)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(op)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Data provides a mock function with given fields:
func (_m *PlaybackHandle) Data() (*ari.PlaybackData, error) {
	ret := _m.Called()

	var r0 *ari.PlaybackData
	if rf, ok := ret.Get(0).(func() *ari.PlaybackData); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.PlaybackData)
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

// ID provides a mock function with given fields:
func (_m *PlaybackHandle) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Match provides a mock function with given fields: e
func (_m *PlaybackHandle) Match(e ari.Event) bool {
	ret := _m.Called(e)

	var r0 bool
	if rf, ok := ret.Get(0).(func(ari.Event) bool); ok {
		r0 = rf(e)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *PlaybackHandle) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: n
func (_m *PlaybackHandle) Subscribe(n ...string) ari.Subscription {
	_va := make([]interface{}, len(n))
	for _i := range n {
		_va[_i] = n[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 ari.Subscription
	if rf, ok := ret.Get(0).(func(...string) ari.Subscription); ok {
		r0 = rf(n...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ari.Subscription)
		}
	}

	return r0
}
