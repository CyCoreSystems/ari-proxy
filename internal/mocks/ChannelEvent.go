package mocks

import mock "github.com/stretchr/testify/mock"

// ChannelEvent is an autogenerated mock type for the ChannelEvent type
type ChannelEvent struct {
	mock.Mock
}

// GetChannelIDs provides a mock function with given fields:
func (_m *ChannelEvent) GetChannelIDs() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
