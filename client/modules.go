package client

/*
type natsModules struct {
	conn *Conn
}

func (m *natsModules) Get(name string) *ari.ModuleHandle {
	return ari.NewModuleHandle(name, m)
}

func (m *natsModules) List() (mx []*ari.ModuleHandle, err error) {
	var modules []string
	err = m.conn.ReadRequest("ari.modules.all", "", nil, &modules)
	for _, mh := range modules {
		mx = append(mx, m.Get(mh))
	}
	return
}

func (m *natsModules) Reload(name string) (err error) {
	err = m.conn.StandardRequest("ari.modules.reload", name, nil, nil)
	return
}

func (m *natsModules) Unload(name string) (err error) {
	err = m.conn.StandardRequest("ari.modules.unload", name, nil, nil)
	return
}

func (m *natsModules) Load(name string) (err error) {
	err = m.conn.StandardRequest("ari.modules.load", name, nil, nil)
	return
}

func (m *natsModules) Data(name string) (md ari.ModuleData, err error) {
	err = m.conn.ReadRequest("ari.modules.data", name, nil, &md)
	return
}
*/
