package integration

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/arimocks"
	tmock "github.com/stretchr/testify/mock"
)

type mock struct {
	Bus    *arimocks.Bus
	Client *arimocks.Client

	Application   *arimocks.Application
	Asterisk      *arimocks.Asterisk
	Bridge        *arimocks.Bridge
	Channel       *arimocks.Channel
	DeviceState   *arimocks.DeviceState
	Endpoint      *arimocks.Endpoint
	LiveRecording *arimocks.LiveRecording
	Logging       *arimocks.Logging
	Mailbox       *arimocks.Mailbox
	Modules       *arimocks.Modules
	Playback      *arimocks.Playback
	Sound         *arimocks.Sound

	AllSub          *arimocks.Subscription
	AllEventChannel <-chan ari.Event

	Shutdown func()
}

func standardMock() *mock {
	m := &mock{}

	m.Bus = &arimocks.Bus{}
	m.Client = &arimocks.Client{}

	m.Asterisk = &arimocks.Asterisk{}
	m.Application = &arimocks.Application{}
	m.Bridge = &arimocks.Bridge{}
	m.Channel = &arimocks.Channel{}
	m.DeviceState = &arimocks.DeviceState{}
	m.Endpoint = &arimocks.Endpoint{}
	m.LiveRecording = &arimocks.LiveRecording{}
	m.Logging = &arimocks.Logging{}
	m.Mailbox = &arimocks.Mailbox{}
	m.Modules = &arimocks.Modules{}
	m.Playback = &arimocks.Playback{}
	m.Sound = &arimocks.Sound{}

	m.AllSub = &arimocks.Subscription{}

	eventCh := make(<-chan ari.Event)

	m.AllSub.On("Cancel").Return(nil)
	m.AllSub.On("Events").Return(eventCh)
	m.Bus.On("Subscribe", tmock.Anything, "all").Return(m.AllSub).Times(1)

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
	m.Client.On("Playback").Return(m.Playback)
	m.Client.On("Sound").Return(m.Sound)

	m.Asterisk.On("Info", (*ari.Key)(nil)).Return(&ari.AsteriskInfo{
		SystemInfo: ari.SystemInfo{
			EntityID: "1",
		},
	}, nil).Times(1) // ensure that downstream tests of Info do not struggle to test additional cases

	return m
}
