package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type endpoint struct {
	c *Client
}

func (e *endpoint) Data(key *ari.Key) (*ari.EndpointData, error) {
	data, err := e.c.dataRequest(&proxy.Request{
		Kind: "EndpointData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Endpoint, nil
}

func (e *endpoint) Get(key *ari.Key) *ari.EndpointHandle {
	k, err := e.c.getRequest(&proxy.Request{
		Kind: "EndpointGet",
		Key:  key,
	})
	if err != nil {
		e.c.log.Warn("failed to get endpoint for handle", "error", err)
		return ari.NewEndpointHandle(key, e)
	}
	return ari.NewEndpointHandle(k, e)
}

func (e *endpoint) List(filter *ari.Key) ([]*ari.Key, error) {
	return e.c.listRequest(&proxy.Request{
		Kind: "EndpointList",
		Key:  filter,
	})
}

func (e *endpoint) ListByTech(tech string, filter *ari.Key) ([]*ari.Key, error) {
	return e.c.listRequest(&proxy.Request{
		Kind: "EndpointListByTech",
		Key:  filter,
		EndpointListByTech: &proxy.EndpointListByTech{
			Tech: tech,
		},
	})
}
