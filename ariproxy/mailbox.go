package ariproxy

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari-proxy/session"
)

func (ins *Instance) mailbox() {
	ins.subscribe("ari.mailboxes.all", func(msg *session.Message, reply Reply) {
		mx, err := ins.upstream.Mailbox.List()
		if err != nil {
			reply(nil, err)
			return
		}

		var mailboxes []string
		for _, m := range mx {
			mailboxes = append(mailboxes, m.ID())
		}

		reply(mailboxes, nil)
	})

	ins.subscribe("ari.mailboxes.data", func(msg *session.Message, reply Reply) {
		name := msg.Object
		dd, err := ins.upstream.Mailbox.Data(name)
		reply(dd, err)
	})

	ins.subscribe("ari.mailboxes.update", func(msg *session.Message, reply Reply) {
		name := msg.Object

		type req struct {
			Old int `json:"old"`
			New int `json:"new"`
		}

		var request req
		if err := json.Unmarshal(msg.Payload, &request); err != nil {
			reply(nil, &decodingError{msg.Command, err})
			return
		}

		err := ins.upstream.Mailbox.Update(name, request.Old, request.New)
		reply(nil, err)
	})

	ins.subscribe("ari.mailboxes.delete", func(msg *session.Message, reply Reply) {
		name := msg.Object
		err := ins.upstream.Mailbox.Delete(name)
		reply(nil, err)
	})

}
