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

	req := proxy.Request{
		SoundList: &proxy.SoundList{
			Filters: filters,
		},
	}
	var resp proxy.Response
	err = s.c.nc.Request(proxy.GetSubject(s.c.prefix, s.c.appName, ""), &req, &resp, s.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
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
	req := proxy.Request{
		SoundData: &proxy.SoundData{
			Name: name,
		},
	}
	var resp proxy.DataResponse
	err = s.c.nc.Request(proxy.GetSubject(s.c.prefix, s.c.appName, ""), &req, &resp, s.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	sd = resp.SoundData
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
