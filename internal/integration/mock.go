package integration

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/internal/mocks"
)

type mock struct {
	Bus         *mocks.Bus
	Client      *mocks.Client
	Asterisk    *mocks.Asterisk
	Application *mocks.Application

	AllSub          *mocks.Subscription
	AllEventChannel <-chan ari.Event
}

func standardMock() *mock {
	m := &mock{}
	m.Bus = &mocks.Bus{}
	m.Client = &mocks.Client{}
	m.Asterisk = &mocks.Asterisk{}
	m.Application = &mocks.Application{}
	m.AllSub = &mocks.Subscription{}

	eventCh := make(<-chan ari.Event)

	m.AllSub.On("Cancel").Return(nil)
	m.AllSub.On("Events").Return(eventCh)
	m.Bus.On("Subscribe", "all").Return(m.AllSub)

	m.Client.On("Bus").Return(m.Bus)

	m.Client.On("ApplicationName").Return("asdf")
	m.Client.On("Asterisk").Return(m.Asterisk)
	m.Client.On("Application").Return(m.Application)

	m.Asterisk.On("Info", "").Return(&ari.AsteriskInfo{
		SystemInfo: ari.SystemInfo{
			EntityID: "1",
		},
	}, nil)

	return m
}
