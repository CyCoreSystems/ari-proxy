package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type application struct {
	c *Client
}

func (a *application) List(filter *ari.Key) ([]*ari.Key, error) {
	return a.c.listRequest(&proxy.Request{
		Kind: "ApplicationList",
		Key:  filter,
	})
}

func (a *application) Data(key *ari.Key) (*ari.ApplicationData, error) {
	ret, err := a.c.dataRequest(&proxy.Request{
		Kind: "ApplicationData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return ret.Application, nil
}

func (a *application) Get(key *ari.Key) *ari.ApplicationHandle {
	k, err := a.c.getRequest(&proxy.Request{
		Kind: "ApplicationGet",
		Key:  key,
	})
	if err != nil {
		a.c.log.Warn("failed to make data request for application", "error", err)
		return ari.NewApplicationHandle(key, a)
	}
	return ari.NewApplicationHandle(k, a)
}

func (a *application) Subscribe(key *ari.Key, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		Kind: "ApplicationSubscribe",
		Key:  key,
		ApplicationSubscribe: &proxy.ApplicationSubscribe{
			EventSource: eventSource,
		},
	})
}

func (a *application) Unsubscribe(key *ari.Key, eventSource string) (err error) {
	return a.c.commandRequest(&proxy.Request{
		Kind: "ApplicationUnsubscribe",
		Key:  key,
		ApplicationSubscribe: &proxy.ApplicationSubscribe{
			EventSource: eventSource,
		},
	})
}
