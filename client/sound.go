package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type sound struct {
	c *Client
}

func (s *sound) List(filters map[string]string) (ret []ari.SoundHandle, err error) {

	if filters == nil {
		filters = make(map[string]string)
	}

	el, err := s.c.listRequest(&proxy.Request{
		SoundList: &proxy.SoundList{
			Filters: filters,
		},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
		ret = append(ret, s.Get(i.ID))
	}
	return
}

func (s *sound) Get(name string) ari.SoundHandle {
	return &soundHandle{
		id: name,
		s:  s,
	}
}

func (s *sound) Data(name string) (sd *ari.SoundData, err error) {
	data, err := s.c.dataRequest(&proxy.Request{
		SoundData: &proxy.SoundData{
			Name: name,
		},
	})
	if err != nil {
		return
	}
	sd = data.Sound
	return
}

type soundHandle struct {
	id string
	s  *sound
}

func (s *soundHandle) Data() (*ari.SoundData, error) {
	return s.s.Data(s.id)
}

func (s *soundHandle) ID() string {
	return s.id
}
