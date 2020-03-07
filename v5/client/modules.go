package client

import (
	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type modules struct {
	c *Client
}

func (m *modules) Data(key *ari.Key) (*ari.ModuleData, error) {
	data, err := m.c.dataRequest(&proxy.Request{
		Kind: "AsteriskModuleData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Module, nil
}

func (m *modules) Get(key *ari.Key) *ari.ModuleHandle {
	k, err := m.c.getRequest(&proxy.Request{
		Kind: "AsteriskModuleGet",
		Key:  key,
	})
	if err != nil {
		m.c.log.Warn("failed to get module for handle", "error", err)
		return ari.NewModuleHandle(key, m)
	}
	return ari.NewModuleHandle(k, m)
}

func (m *modules) List(filter *ari.Key) ([]*ari.Key, error) {
	return m.c.listRequest(&proxy.Request{
		Kind: "AsteriskModuleList",
		Key:  filter,
	})
}

func (m *modules) Load(key *ari.Key) error {
	return m.c.commandRequest(&proxy.Request{
		Kind: "AsteriskModuleLoad",
		Key:  key,
	})
}

func (m *modules) Reload(key *ari.Key) error {
	return m.c.commandRequest(&proxy.Request{
		Kind: "AsteriskModuleReload",
		Key:  key,
	})
}

func (m *modules) Unload(key *ari.Key) error {
	return m.c.commandRequest(&proxy.Request{
		Kind: "AsteriskModuleUnload",
		Key:  key,
	})
}
