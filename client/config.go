package client

import (
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type config struct {
	c *Client
}

func (c *config) Get(key *ari.Key) *ari.ConfigHandle {
	return ari.NewConfigHandle(key, c)
}

func (c *config) Data(key *ari.Key) (*ari.ConfigData, error) {
	data, err := c.c.dataRequest(&proxy.Request{
		Kind: "AsteriskConfigData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Config, nil
}

func (c *config) Update(key *ari.Key, tuples []ari.ConfigTuple) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "AsteriskConfigUpdate",
		Key:  key,
		AsteriskConfig: &proxy.AsteriskConfig{
			Tuples: tuples,
		},
	})
}

func (c *config) Delete(key *ari.Key) error {
	return c.c.commandRequest(&proxy.Request{
		Kind: "AsteriskConfigDelete",
		Key:  key,
	})
}
