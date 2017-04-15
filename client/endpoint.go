package client

import (
	"fmt"
	"strings"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type endpoint struct {
	c *Client
}

func (e *endpoint) Data(tech string, resource string) (ed *ari.EndpointData, err error) {
	d, err := e.c.dataRequest(&proxy.Request{
		EndpointData: &proxy.EndpointData{
			Resource: resource,
			Tech:     tech,
		},
	})
	if err != nil {
		return
	}
	ed = d.Endpoint
	return
}

func (e *endpoint) Get(tech string, resource string) ari.EndpointHandle {
	return &endpointHandle{
		tech:     tech,
		resource: resource,
		e:        e,
	}
}

func (e *endpoint) List() (ret []ari.EndpointHandle, err error) {
	el, err := e.c.listRequest(&proxy.Request{
		EndpointList: &proxy.EndpointList{},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
		items := strings.Split(i.ID, "/")
		ret = append(ret, e.Get(items[0], items[1]))
	}
	return
}

func (e *endpoint) ListByTech(tech string) (ret []ari.EndpointHandle, err error) {
	el, err := e.c.listRequest(&proxy.Request{
		EndpointListByTech: &proxy.EndpointListByTech{
			Tech: tech,
		},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
		items := strings.Split(i.ID, "/")
		ret = append(ret, e.Get(items[0], items[1]))
	}
	return
}

type endpointHandle struct {
	tech     string
	resource string
	e        *endpoint
}

func (e *endpointHandle) Data() (*ari.EndpointData, error) {
	return e.e.Data(e.tech, e.resource)
}

func (e *endpointHandle) ID() string {
	return fmt.Sprintf("%s/%s", e.tech, e.resource)
}

func (e *endpointHandle) Match(ev ari.Event) (ok bool) {
	v, ok := ev.(ari.EndpointEvent)
	if !ok {
		return false
	}
	for _, i := range v.GetEndpointIDs() {
		if i == e.ID() {
			return true
		}
	}
	return false
}
