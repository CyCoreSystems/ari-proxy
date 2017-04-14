package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type deviceState struct {
	c *Client
}

func (ds *deviceState) Get(name string) ari.DeviceStateHandle {
	//return ari.NewDeviceStateHandle(name, ds)
	return nil
}

func (ds *deviceState) List() (dx []ari.DeviceStateHandle, err error) {
	req := proxy.Request{
		DeviceStateList: &proxy.DeviceStateList{},
	}
	var resp proxy.Response
	err = ds.c.nc.Request(proxy.GetSubject(ds.c.prefix, ds.c.appName, ""), &req, &resp, ds.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		dx = append(dx, ds.Get(i.ID))
	}
	return
}

func (ds *deviceState) Data(name string) (d *ari.DeviceStateData, err error) {
	req := proxy.Request{
		DeviceStateData: &proxy.DeviceStateData{
			ID: name,
		},
	}
	var resp proxy.DataResponse
	err = ds.c.nc.Request(proxy.GetSubject(ds.c.prefix, ds.c.appName, ""), &req, &resp, ds.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	d = resp.DeviceStateData
	return
}

func (ds *deviceState) Update(name string, state string) (err error) {
	req := proxy.Request{
		DeviceStateUpdate: &proxy.DeviceStateUpdate{
			ID:    name,
			State: state,
		},
	}
	var resp proxy.Response
	err = ds.c.nc.Request(proxy.CommandSubject(ds.c.prefix, ds.c.appName, ""), &req, &resp, ds.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (ds *deviceState) Delete(name string) (err error) {
	req := proxy.Request{
		DeviceStateDelete: &proxy.DeviceStateDelete{
			ID: name,
		},
	}
	var resp proxy.Response
	err = ds.c.nc.Request(proxy.CommandSubject(ds.c.prefix, ds.c.appName, ""), &req, &resp, ds.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}
