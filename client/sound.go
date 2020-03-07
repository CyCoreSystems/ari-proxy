package client

import (
	"github.com/CyCoreSystems/ari-proxy/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type sound struct {
	c *Client
}

func (s *sound) List(filters map[string]string, keyFilter *ari.Key) ([]*ari.Key, error) {
	return s.c.listRequest(&proxy.Request{
		Kind: "SoundList",
		Key:  keyFilter,
		SoundList: &proxy.SoundList{
			Filters: filters,
		},
	})
}

func (s *sound) Data(key *ari.Key) (*ari.SoundData, error) {
	data, err := s.c.dataRequest(&proxy.Request{
		Kind: "SoundData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Sound, nil
}
