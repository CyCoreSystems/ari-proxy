package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type asterisk struct {
	c *Client
}

func (a *asterisk) Config() ari.Config {
	return &config{a.c}
}

type asteriskVariables struct {
	c *Client
}

func (a *asterisk) Logging() ari.Logging {
	return nil
}

func (a *asterisk) Modules() ari.Modules {
	return nil
}

func (a *asterisk) ReloadModule(name string) (err error) {
	err = a.Modules().Reload(name)
	return
}

func (a *asterisk) Info(only string) (*ari.AsteriskInfo, error) {
	resp, err := a.c.dataRequest(&proxy.Request{
		AsteriskInfo: &proxy.AsteriskInfo{},
	})
	if err != nil {
		return nil, err
	}
	return resp.Asterisk, nil
}

func (a *asterisk) Variables() ari.Variables {
	return &asteriskVariables{a.c}
}

func (a *asteriskVariables) Get(key string) (ret string, err error) {
	data, err := a.c.dataRequest(&proxy.Request{
		AsteriskVariables: &proxy.AsteriskVariables{
			Name: key,
			Get:  &proxy.VariablesGet{},
		},
	})
	if err != nil {
		return "", err
	}
	return data.Variable, err
}

func (a *asteriskVariables) Set(key string, val string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		AsteriskVariables: &proxy.AsteriskVariables{
			Name: key,
			Set: &proxy.VariablesSet{
				Value: val,
			},
		},
	})
}
