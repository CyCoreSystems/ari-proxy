package client

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type config struct {
	c *Client
}

func (c *config) Get(configClass string, objectType string, id string) ari.ConfigHandle {
	return &configHandle{
		configClass: configClass,
		objectType:  objectType,
		id:          id,
		c:           c,
	}
}

func (c *config) Data(configClass string, objectType string, id string) (cd *ari.ConfigData, err error) {
	req := proxy.Request{
		AsteriskConfig: &proxy.AsteriskConfig{
			ConfigClass: configClass,
			ObjectType:  objectType,
			ID:          id,
			Data:        &proxy.AsteriskConfigData{},
		},
	}
	var resp proxy.DataResponse
	err = c.c.nc.Request(proxy.GetSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	cd = resp.ConfigData
	return
}

func (c *config) Update(configClass string, objectType string, id string, tuples []ari.ConfigTuple) (err error) {
	req := proxy.Request{
		AsteriskConfig: &proxy.AsteriskConfig{
			ConfigClass: configClass,
			ObjectType:  objectType,
			ID:          id,
			Update: &proxy.AsteriskConfigUpdate{
				Tuples: tuples,
			},
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.GetSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (c *config) Delete(configClass string, objectType string, id string) (err error) {
	req := proxy.Request{
		AsteriskConfig: &proxy.AsteriskConfig{
			ConfigClass: configClass,
			ObjectType:  objectType,
			ID:          id,
			Delete:      &proxy.AsteriskConfigDelete{},
		},
	}
	var resp proxy.Response
	err = c.c.nc.Request(proxy.GetSubject(c.c.prefix, c.c.appName, ""), &req, &resp, c.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

type configHandle struct {
	configClass string
	objectType  string
	id          string
	c           *config
}

func (c *configHandle) Data() (*ari.ConfigData, error) {
	return c.c.Data(c.configClass, c.objectType, c.id)
}

func (c *configHandle) Delete() error {
	return c.c.Delete(c.configClass, c.objectType, c.id)
}

func (c *configHandle) ID() string {
	return fmt.Sprintf("%s/%s/%s", c.configClass, c.objectType, c.id)
}

func (c *configHandle) Update(c1 []ari.ConfigTuple) error {
	return c.c.Update(c.configClass, c.objectType, c.id, c1)
}
