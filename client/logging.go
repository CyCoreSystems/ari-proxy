package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type logging struct {
	c *Client
}

func (l *logging) Create(name string, level string) (err error) {
	req := proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Create: &proxy.AsteriskLoggingCreate{
				Config: level,
				ID:     name,
			},
		},
	}
	var resp proxy.DataResponse
	err = l.c.nc.Request(proxy.CreateSubject(l.c.prefix, l.c.appName, ""), &req, &resp, l.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (l *logging) List() (ld []ari.LogData, err error) {
	req := proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			List: &proxy.AsteriskLoggingList{},
		},
	}
	var resp proxy.DataResponse
	err = l.c.nc.Request(proxy.GetSubject(l.c.prefix, l.c.appName, ""), &req, &resp, l.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	ld = resp.LogDataList
	return
}

func (l *logging) Rotate(name string) (err error) {
	req := proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Rotate: &proxy.AsteriskLoggingRotate{
				ID: name,
			},
		},
	}
	var resp proxy.DataResponse
	err = l.c.nc.Request(proxy.CommandSubject(l.c.prefix, l.c.appName, ""), &req, &resp, l.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (l *logging) Delete(name string) (err error) {
	req := proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Delete: &proxy.AsteriskLoggingDelete{
				ID: name,
			},
		},
	}
	var resp proxy.DataResponse
	err = l.c.nc.Request(proxy.CommandSubject(l.c.prefix, l.c.appName, ""), &req, &resp, l.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}
