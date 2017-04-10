package ariproxy

import (
	"encoding/json"
	"errors"

	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
)

func (ins *Instance) subscribe(endpoint string, h Handler2) {
	ins.log.Debug("Registering command", "endpoint", endpoint)
	ins.dispatcherLock.Lock()
	defer ins.dispatcherLock.Unlock()
	ins.dispatcher[endpoint] = h
}

func (ins *Instance) commands() {

	endpoint := "ari.commands.dialog." + ins.Dialog.ID
	ins.log.Debug("Subscribing on nats endpoint", "endpoint", endpoint)

	cb := func(m *nats.Msg) {

		reply := m.Reply
		data := m.Data

		ins.conn.Publish(reply+".ok", []byte(""))

		var msg session.Message

		if err := json.Unmarshal(data, &msg); err != nil {
			respMap := client.ErrorToMap(err, "")
			resp, _ := json.Marshal(respMap)
			if err2 := ins.conn.Publish(reply+".err", resp); err2 != nil {
				ins.log.Error("Error sending error reply", "reply", reply, "error", err2, "body", err)
			}
			return
		}

		ins.dispatcherLock.RLock()

		// find using full name
		h, _ := ins.dispatcher[msg.Command]

		ins.dispatcherLock.RUnlock()

		if h == nil {
			err := errors.New("No handler found for " + msg.Command)
			respMap := client.ErrorToMap(err, "")
			resp, _ := json.Marshal(respMap)
			if err2 := ins.conn.Publish(reply+".err", resp); err2 != nil {
				ins.log.Error("Error sending error reply", "reply", reply, "error", err2, "body", err)
			}
			return
		}

		h(&msg, func(i interface{}, err error) {

			if err != nil {
				respMap := client.ErrorToMap(err, "")
				resp, _ := json.Marshal(respMap)
				if err2 := ins.conn.Publish(reply+".err", resp); err2 != nil {
					ins.log.Error("Error sending error reply", "reply", reply, "error", err2, "body", err)
				}
				return
			}

			resp := []byte("{}")
			if i != nil {
				resp, err = json.Marshal(i)
				if err != nil {
					ins.log.Error("Error building response reply", "reply", reply, "error", err)
					return
				}
			}

			if err = ins.conn.Publish(reply+".resp", resp); err != nil {
				ins.log.Error("Error sending response reply", "reply", reply, "error", err, "body", resp)
			}
		})
	}

	sub, err := ins.conn.Subscribe(endpoint, func(m *nats.Msg) { go cb(m) })
	if err != nil {
		ins.log.Error("Error starting subscription", "endpoint", endpoint, "error", err)
		ins.Stop()
		return
	}

	go func() {
		<-ins.ctx.Done()
		sub.Unsubscribe()
	}()
}
