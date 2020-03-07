package client

import (
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari/v5"
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
	return &logging{a.c}
}

func (a *asterisk) Modules() ari.Modules {
	return &modules{a.c}
}

func (a *asterisk) Info(key *ari.Key) (*ari.AsteriskInfo, error) {
	resp, err := a.c.dataRequest(&proxy.Request{
		Kind: "AsteriskInfo",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Asterisk, nil
}

func (a *asterisk) Variables() ari.AsteriskVariables {
	return &asteriskVariables{a.c}
}

func (a *asteriskVariables) Get(key *ari.Key) (ret string, err error) {
	data, err := a.c.dataRequest(&proxy.Request{
		Kind: "AsteriskVariableGet",
		Key:  key,
	})
	if err != nil {
		return "", err
	}
	return data.Variable, err
}

func (a *asteriskVariables) Set(key *ari.Key, val string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		Kind: "AsteriskVariableSet",
		Key:  key,
		AsteriskVariableSet: &proxy.AsteriskVariableSet{
			Value: val,
		},
	})
}
