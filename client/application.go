package client

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type application struct {
	c *Client
}

func (a *application) List() ([]ari.ApplicationHandle, error) {
	panic("not implemented")
}

func (a *application) Get(name string) ari.ApplicationHandle {
	panic("not implemented")
}

func (a *application) Data(name string) (*ari.ApplicationData, error) {
	panic("not implemented")
}

func (a *application) Subscribe(name string, eventSource string) error {
	panic("not implemented")
}

func (a *application) Unsubscribe(name string, eventSource string) error {
	panic("not implemented")
}

func (a *application) List() (ax []*ari.ApplicationHandle, err error) {
	req := proxy.Request{
		ApplicationList: &proxy.ApplicationList{},
	}
	var resp proxy.EntityList
	err = a.nats.Request(fmt.Sprintf("%sget.%s", a.prefix, a.app), &req, &resp, a.opts.RequestTimeout)
	if err != nil {
		return
	}

	for _, i := range resp.List {
		ret = append(ret, ari.NewApplicationHandle(i.ID(), a.app))
	}
	return
}

func (a *application) Get(name string) *ari.ApplicationHandle {
	panic("not implemented")
}

func (a *application) Data(name string) (ari.ApplicationData, error) {
	panic("not implemented")
}

func (a *application) Subscribe(name string, eventSource string) error {
	panic("not implemented")
}

func (a *application) Unsubscribe(name string, eventSource string) error {
	panic("not implemented")
}

func (a *natsApplication) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

func (a *natsApplication) List() (ax []*ari.ApplicationHandle, err error) {
	var apps []string
	err = a.conn.ReadRequest("ari.applications.all", "", nil, &apps)
	for _, app := range apps {
		ax = append(ax, a.Get(app))
	}
	return
}

func (a *natsApplication) Data(name string) (d ari.ApplicationData, err error) {
	err = a.conn.ReadRequest("ari.applications.data", name, nil, &d)
	return
}

func (a *natsApplication) Subscribe(name string, eventSource string) (err error) {
	err = a.conn.StandardRequest("ari.applications.subscribe", name, eventSource, nil)
	return
}

func (a *natsApplication) Unsubscribe(name string, eventSource string) (err error) {
	err = a.conn.StandardRequest("ari.applications.unsubscribe", name, eventSource, nil)
	return
}
