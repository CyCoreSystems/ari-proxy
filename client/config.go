package client

/*
type natsConfig struct {
	conn *Conn
}

func (c *natsConfig) Get(configClass string, objectType string, id string) *ari.ConfigHandle {
	return ari.NewConfigHandle(configClass, objectType, id, c)
}

func (c *natsConfig) Data(configClass string, objectType string, id string) (cd ari.ConfigData, err error) {
	cd.ID = id
	cd.Type = objectType
	cd.Class = configClass
	err = c.conn.ReadRequest("ari.asterisk.config.data", configClass+"."+objectType+"."+id, nil, &cd.Fields)
	return
}

func (c *natsConfig) Update(configClass string, objectType string, id string, tuples []ari.ConfigTuple) (err error) {
	err = c.conn.StandardRequest("ari.asterisk.config.update", configClass+"."+objectType+"."+id, &tuples, nil)
	return
}

func (c *natsConfig) Delete(configClass string, objectType string, id string) (err error) {
	err = c.conn.StandardRequest("ari.asterisk.config.delete", configClass+"."+objectType+"."+id, nil, nil)
	return
}
*/
