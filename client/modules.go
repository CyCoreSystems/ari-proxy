package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type modules struct {
	c *Client
}

func (m *modules) Get(name string) ari.ModuleHandle {
	return &moduleHandle{
		m:  m,
		id: name,
	}
}

func (m *modules) List() (ret []ari.ModuleHandle, err error) {
	req := proxy.Request{
		AsteriskModules: &proxy.AsteriskModules{
			List: &proxy.AsteriskModulesList{},
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		ret = append(ret, m.Get(i.ID))
	}
	return
}

func (m *modules) Reload(name string) (err error) {
	req := proxy.Request{
		AsteriskModules: &proxy.AsteriskModules{
			Reload: &proxy.AsteriskModulesReload{
				Name: name,
			},
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (m *modules) Unload(name string) (err error) {
	req := proxy.Request{
		AsteriskModules: &proxy.AsteriskModules{
			Unload: &proxy.AsteriskModulesUnload{
				Name: name,
			},
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (m *modules) Load(name string) (err error) {
	req := proxy.Request{
		AsteriskModules: &proxy.AsteriskModules{
			Load: &proxy.AsteriskModulesLoad{
				Name: name,
			},
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (m *modules) Data(name string) (md *ari.ModuleData, err error) {
	req := proxy.Request{
		AsteriskModules: &proxy.AsteriskModules{
			Data: &proxy.AsteriskModulesData{
				Name: name,
			},
		},
	}
	var resp proxy.DataResponse
	err = m.c.nc.Request(proxy.GetSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	md = resp.ModuleData
	return
}

type moduleHandle struct {
	id string
	m  *modules
}

func (m *moduleHandle) Data() (*ari.ModuleData, error) {
	return m.m.Data(m.id)
}

func (m *moduleHandle) ID() string {
	return m.id
}

func (m *moduleHandle) Load() error {
	return m.m.Load(m.id)
}

func (m *moduleHandle) Reload() error {
	return m.m.Reload(m.id)
}

func (m *moduleHandle) Unload() error {
	return m.m.Unload(m.id)
}
