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
	el, err := ds.c.listRequest(&proxy.Request{
		DeviceStateList: &proxy.DeviceStateList{},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
		dx = append(dx, ds.Get(i.ID))
	}
	return
}

func (ds *deviceState) Data(name string) (d *ari.DeviceStateData, err error) {
	nd, err := ds.c.dataRequest(&proxy.Request{
		DeviceStateData: &proxy.DeviceStateData{
			ID: name,
		},
	})
	if err != nil {
		return
	}
	d = nd.DeviceState
	return
}

func (ds *deviceState) Update(name string, state string) (err error) {
	err = ds.c.commandRequest(&proxy.Request{
		DeviceStateUpdate: &proxy.DeviceStateUpdate{
			ID:    name,
			State: state,
		},
	})
	return
}

func (ds *deviceState) Delete(name string) (err error) {
	err = ds.c.commandRequest(&proxy.Request{
		DeviceStateDelete: &proxy.DeviceStateDelete{
			ID: name,
		},
	})
	return
}
