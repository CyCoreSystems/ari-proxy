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
	req := proxy.Request{
		EndpointData: &proxy.EndpointData{
			Resource: resource,
			Tech:     tech,
		},
	}
	var resp proxy.DataResponse
	err = e.c.nc.Request(proxy.GetSubject(e.c.prefix, e.c.appName, ""), &req, &resp, e.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	ed = resp.EndpointData
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
	req := proxy.Request{
		EndpointList: &proxy.EndpointList{},
	}
	var resp proxy.Response
	err = e.c.nc.Request(proxy.GetSubject(e.c.prefix, e.c.appName, ""), &req, &resp, e.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		items := strings.Split(i.ID, "/")
		ret = append(ret, e.Get(items[0], items[1]))
	}
	return
}

func (e *endpoint) ListByTech(tech string) (ret []ari.EndpointHandle, err error) {
	req := proxy.Request{
		EndpointListByTech: &proxy.EndpointListByTech{
			Tech: tech,
		},
	}
	var resp proxy.Response
	err = e.c.nc.Request(proxy.GetSubject(e.c.prefix, e.c.appName, ""), &req, &resp, e.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
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

func (e *endpointHandle) Match(e1 ari.Event) (ok bool) {
	return
}