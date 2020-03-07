package client

import (
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type logging struct {
	c *Client
}

func (l *logging) Create(key *ari.Key, levels string) (*ari.LogHandle, error) {
	k, err := l.c.createRequest(&proxy.Request{
		Kind: "AsteriskLoggingCreate",
		Key:  key,
		AsteriskLoggingChannel: &proxy.AsteriskLoggingChannel{
			Levels: levels,
		},
	})
	if err != nil {
		return nil, err
	}
	return ari.NewLogHandle(k, l), nil
}

func (l *logging) Data(key *ari.Key) (*ari.LogData, error) {
	data, err := l.c.dataRequest(&proxy.Request{
		Kind: "AsteriskLoggingData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Log, nil
}

func (l *logging) Get(key *ari.Key) *ari.LogHandle {
	k, err := l.c.getRequest(&proxy.Request{
		Kind: "AsteriskLoggingGet",
		Key:  key,
	})
	if err != nil {
		l.c.log.Warn("failed to get logging key for handle", "error", err)
		return ari.NewLogHandle(key, l)
	}
	return ari.NewLogHandle(k, l)
}

func (l *logging) List(filter *ari.Key) ([]*ari.Key, error) {
	return l.c.listRequest(&proxy.Request{
		Kind: "AsteriskLoggingList",
		Key:  filter,
	})
}

func (l *logging) Rotate(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "AsteriskLoggingRotate",
		Key:  key,
	})
}

func (l *logging) Delete(key *ari.Key) error {
	return l.c.commandRequest(&proxy.Request{
		Kind: "AsteriskLoggingDelete",
		Key:  key,
	})
}
