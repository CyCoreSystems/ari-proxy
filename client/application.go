package client

import "github.com/CyCoreSystems/ari"

type natsApplication struct {
	conn *Conn
}

func (a *natsApplication) Get(name string) *ari.ApplicationHandle {
	return ari.NewApplicationHandle(name, a)
}

func (a *natsApplication) List() (ax []*ari.ApplicationHandle, err error) {
	var apps []string
	err = a.conn.ReadRequest("ari.applications.all", "", nil, &apps)
	for _, app := range apps {
		ax = append(ax, a.Get(app))
	}
	return
}

func (a *natsApplication) Data(name string) (d ari.ApplicationData, err error) {
	err = a.conn.ReadRequest("ari.applications.data", name, nil, &d)
	return
}

func (a *natsApplication) Subscribe(name string, eventSource string) (err error) {
	err = a.conn.StandardRequest("ari.applications.subscribe", name, eventSource, nil)
	return
}

func (a *natsApplication) Unsubscribe(name string, eventSource string) (err error) {
	err = a.conn.StandardRequest("ari.applications.unsubscribe", name, eventSource, nil)
	return
}
