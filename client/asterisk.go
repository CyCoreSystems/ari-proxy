package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type asterisk struct {
	c *Client
}

func (a *asterisk) Config() ari.Config {
	return nil
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

func (a *asterisk) Info(only string) (ai *ari.AsteriskInfo, err error) {
	req := proxy.Request{
		AsteriskInfo: &proxy.AsteriskInfo{},
	}
	var resp proxy.DataResponse
	err = a.c.nc.Request(proxy.GetSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	ai = resp.AsteriskInfo
	return
}

func (a *asterisk) Variables() ari.Variables {
	return &asteriskVariables{a.c}
}

func (a *asteriskVariables) Get(variable string) (ret string, err error) {
	req := proxy.Request{
		AsteriskVariables: &proxy.AsteriskVariables{
			Name: variable,
			Get:  &proxy.VariablesGet{},
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
	ret = resp.Variable
	return
}

func (a *asteriskVariables) Set(variable string, value string) (err error) {
	req := proxy.Request{
		AsteriskVariables: &proxy.AsteriskVariables{
			Name: variable,
			Set: &proxy.VariablesSet{
				Value: value,
			},
		},
	}
	var resp proxy.Response
	err = a.c.nc.Request(proxy.CommandSubject(a.c.prefix, a.c.appName, ""), &req, &resp, a.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}
