package integration

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

type mock struct {
	Bus    *mocks.Bus
	Client *mocks.Client

	Application   *mocks.Application
	Asterisk      *mocks.Asterisk
	Bridge        *mocks.Bridge
	Channel       *mocks.Channel
	DeviceState   *mocks.DeviceState
	Endpoint      *mocks.Endpoint
	LiveRecording *mocks.LiveRecording
	Logging       *mocks.Logging
	Mailbox       *mocks.Mailbox
	Modules       *mocks.Modules

	AllSub          *mocks.Subscription
	AllEventChannel <-chan ari.Event

	Shutdown func()
}

func standardMock() *mock {
	m := &mock{}

	m.Bus = &mocks.Bus{}
	m.Client = &mocks.Client{}

	m.Asterisk = &mocks.Asterisk{}
	m.Application = &mocks.Application{}
	m.Bridge = &mocks.Bridge{}
	m.Channel = &mocks.Channel{}
	m.DeviceState = &mocks.DeviceState{}
	m.Endpoint = &mocks.Endpoint{}
	m.LiveRecording = &mocks.LiveRecording{}
	m.Logging = &mocks.Logging{}
	m.Mailbox = &mocks.Mailbox{}
	m.Modules = &mocks.Modules{}

	m.AllSub = &mocks.Subscription{}

	eventCh := make(<-chan ari.Event)

	m.AllSub.On("Cancel").Return(nil)
	m.AllSub.On("Events").Return(eventCh)
	m.Bus.On("Subscribe", "all").Return(m.AllSub).Times(1)

	m.Client.On("Bus").Return(m.Bus)

	m.Client.On("ApplicationName").Return("asdf")
	m.Client.On("Asterisk").Return(m.Asterisk)
	m.Client.On("Application").Return(m.Application)
	m.Client.On("Bridge").Return(m.Bridge)
	m.Client.On("Channel").Return(m.Channel)
	m.Client.On("DeviceState").Return(m.DeviceState)
	m.Client.On("Endpoint").Return(m.Endpoint)
	m.Client.On("LiveRecording").Return(m.LiveRecording)
	m.Asterisk.On("Logging").Return(m.Logging)
	m.Client.On("Mailbox").Return(m.Mailbox)
	m.Asterisk.On("Modules").Return(m.Modules)

	m.Asterisk.On("Info", "").Return(&ari.AsteriskInfo{
		SystemInfo: ari.SystemInfo{
			EntityID: "1",
		},
	}, nil).Times(1) // ensure that downstream tests of Info do not struggle to test additional cases

	return m
}
