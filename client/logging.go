package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type logging struct {
	c *Client
}

func (l *logging) Create(name string, level string) (err error) {
	err = l.c.commandRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Create: &proxy.AsteriskLoggingCreate{
				Config: level,
				ID:     name,
			},
		},
	})
	if err != nil {
		return
	}
	return
}

func (l *logging) List() (ld []ari.LogData, err error) {
	ed, err := l.c.dataRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			List: &proxy.AsteriskLoggingList{},
		},
	})
	if err != nil {
		return
	}

	ld = ed.LogList
	return
}

func (l *logging) Rotate(name string) (err error) {
	err = l.c.commandRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Rotate: &proxy.AsteriskLoggingRotate{
				ID: name,
			},
		},
	})
	return
}

func (l *logging) Delete(name string) (err error) {
	err = l.c.commandRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Delete: &proxy.AsteriskLoggingDelete{
				ID: name,
			},
		},
	})
	return
}
