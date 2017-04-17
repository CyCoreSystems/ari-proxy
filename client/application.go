package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type application struct {
	c *Client
}

func (a *application) List() ([]ari.ApplicationHandle, error) {
	return a.ListMD(nil)
}

// ListMD implements List with metadata
func (a *application) ListMD(m *proxy.Metadata) (ret []ari.ApplicationHandle, err error) {
	list, err := a.c.listRequest(&proxy.Request{
		Metadata:        m,
		ApplicationList: &proxy.ApplicationList{},
	})
	if err != nil {
		return
	}
	for _, i := range list.List {
		ret = append(ret, a.GetMD(i.Metadata, i.ID))
	}
	return
}

func (a *application) Data(name string) (*ari.ApplicationData, error) {
	return a.DataMD(nil, name)
}

func (a *application) DataMD(m *proxy.Metadata, name string) (*ari.ApplicationData, error) {
	app, _, err := a.data(m, name)
	return app, err
}

func (a *application) data(m *proxy.Metadata, name string) (*ari.ApplicationData, *proxy.Metadata, error) {
	ret, err := a.c.dataRequest(&proxy.Request{
		Metadata: m,
		ApplicationData: &proxy.ApplicationData{
			Name: name,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return ret.Application, ret.Metadata, nil
}

func (a *application) Get(id string) (h ari.ApplicationHandle) {
	_, md, err := a.data(nil, id) // obtain metadata by making a data request
	if err != nil {
		a.c.log.Warn("failed to make data request for application", "error", err)
	}
	return a.GetMD(md, id)
}

func (a *application) GetMD(m *proxy.Metadata, id string) (h ari.ApplicationHandle) {
	// make a data request to get the metadata
	return &applicationHandle{id: id, metadata: m, application: a}
}

func (a *application) Subscribe(name string, eventSource string) (err error) {
	return a.SubscribeMD(nil, name, eventSource)
}

func (a *application) SubscribeMD(m *proxy.Metadata, name string, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		Metadata: m,
		ApplicationSubscribe: &proxy.ApplicationSubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	})
}

func (a *application) Unsubscribe(name string, eventSource string) (err error) {
	return a.UnsubscribeMD(nil, name, eventSource)
}

func (a *application) UnsubscribeMD(m *proxy.Metadata, name string, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		Metadata: m,
		ApplicationUnsubscribe: &proxy.ApplicationUnsubscribe{
			EventSource: eventSource,
			Name:        name,
		},
	})
}

type applicationHandle struct {
	id          string
	application *application
	metadata    *proxy.Metadata
}

func (a *applicationHandle) ID() string {
	return a.id
}

func (a *applicationHandle) Data() (ad *ari.ApplicationData, err error) {
	ad, err = a.application.DataMD(a.metadata, a.id)
	return
}

func (a *applicationHandle) Subscribe(eventSource string) (err error) {
	err = a.application.SubscribeMD(a.metadata, a.id, eventSource)
	return
}

func (a *applicationHandle) Unsubscribe(eventSource string) (err error) {
	err = a.application.UnsubscribeMD(a.metadata, a.id, eventSource)
	return
}

func (a *applicationHandle) Match(evt ari.Event) (ok bool) {
	return evt.GetApplication() == a.id
}
