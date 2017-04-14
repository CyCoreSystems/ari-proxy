package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type application struct {
	c *Client
}

func (a *application) List() (ret []ari.ApplicationHandle, err error) {
	list, err := a.c.listRequest(&proxy.Request{
		ApplicationList: &proxy.ApplicationList{},
	})
	if err != nil {
		return
	}
	for _, i := range list.List {
		ret = append(ret, a.Get(i.ID))
	}
	return
}

func (a *application) Data(name string) (*ari.ApplicationData, error) {
	ret, err := a.c.dataRequest(&proxy.Request{
		ApplicationData: &proxy.ApplicationData{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}
	return ret.Application, nil
}

func (a *application) Get(id string) (h ari.ApplicationHandle) {
	return &applicationHandle{id: id, application: a}
}

func (a *application) Subscribe(name string, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		ApplicationSubscribe: &proxy.ApplicationSubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	})
}

func (a *application) Unsubscribe(name string, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		ApplicationUnsubscribe: &proxy.ApplicationUnsubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	})
}

type applicationHandle struct {
	id          string
	application ari.Application
}

func (a *applicationHandle) ID() string {
	return a.id
}

func (a *applicationHandle) Data() (ad *ari.ApplicationData, err error) {
	ad, err = a.application.Data(a.id)
	return
}

func (a *applicationHandle) Subscribe(eventSource string) (err error) {
	err = a.application.Subscribe(a.id, eventSource)
	return
}

func (a *applicationHandle) Unsubscribe(eventSource string) (err error) {
	err = a.application.Unsubscribe(a.id, eventSource)
	return
}

func (a *applicationHandle) Match(evt ari.Event) (ok bool) {
	return evt.GetApplication() == a.id
}
