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

func (l *logging) Data(name string) (*ari.LogData, error) {
	resp, err := l.c.dataRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			Data: &proxy.AsteriskLoggingData{
				ID: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Log, nil
}

func (l *logging) Get(id string) ari.LogHandle {
	return &loggingHandle{
		name: id,
		c:    l,
	}

}

func (l *logging) List() (ld []ari.LogHandle, err error) {
	resp, err := l.c.listRequest(&proxy.Request{
		AsteriskLogging: &proxy.AsteriskLogging{
			List: &proxy.AsteriskLoggingList{},
		},
	})
	if err != nil {
		return
	}

	for _, i := range resp.List {
		ld = append(ld, &loggingHandle{
			metadata: i.Metadata,
			name:     i.ID,
			c:        l,
		})

	}

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

type loggingHandle struct {
	name     string
	c        *logging
	metadata *proxy.Metadata
}

func (l *loggingHandle) ID() string {
	return l.name
}

func (l *loggingHandle) Data() (*ari.LogData, error) {
	return l.c.Data(l.name)
}

func (l *loggingHandle) Rotate() error {
	return l.c.Rotate(l.name)
}

func (l *loggingHandle) Delete() error {
	return l.c.Delete(l.name)
}
