package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type application struct {
	c *Client
}

func (a *application) List() (ret []ari.ApplicationHandle, err error) {
	req := proxy.Request{
		ApplicationList: &proxy.ApplicationList{},
	}
	var resp proxy.EntityList
	err = a.c.nc.Request(proxy.GetSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}

	for _, i := range resp.List {
		ret = append(ret, a.Get(i.ID))
	}
	return
}

func (a *application) Data(name string) (d *ari.ApplicationData, err error) {
	req := proxy.Request{
		ApplicationData: &proxy.ApplicationData{
			Name: name,
		},
	}
	var resp proxy.DataResponse
	err = a.c.nc.Request(proxy.GetSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	d = resp.ApplicationData
	return
}

func (a *application) Get(id string) (h ari.ApplicationHandle) {
	return &applicationHandle{id: id, application: a}
}

func (a *application) Subscribe(name string, eventSource string) (err error) {
	req := proxy.Request{
		ApplicationSubscribe: &proxy.ApplicationSubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	}
	var resp proxy.Response
	err = a.c.nc.Request(proxy.CommandSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}
	err = resp.Err()
	return
}

func (a *application) Unsubscribe(name string, eventSource string) (err error) {
	req := proxy.Request{
		ApplicationUnsubscribe: &proxy.ApplicationUnsubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	}
	var resp proxy.Response
	err = a.c.nc.Request(proxy.CommandSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}
	err = resp.Err()
	return
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
	return
}
