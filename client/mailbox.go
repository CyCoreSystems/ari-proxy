package client

import "github.com/CyCoreSystems/ari"

type natsMailbox struct {
	conn *Conn
}

func (m *natsMailbox) Get(name string) *ari.MailboxHandle {
	return ari.NewMailboxHandle(name, m)
}

func (m *natsMailbox) List() (mx []*ari.MailboxHandle, err error) {
	var boxes []string
	err = m.conn.ReadRequest("ari.mailboxes.all", "", nil, &boxes)
	for _, id := range boxes {
		mx = append(mx, m.Get(id))
	}
	return
}

func (m *natsMailbox) Data(name string) (d ari.MailboxData, err error) {
	err = m.conn.ReadRequest("ari.mailboxes.data", name, nil, &d)
	return
}

// UpdateMailboxRequest is the encoded request for updating the mailbox
type UpdateMailboxRequest struct {
	Old int `json:"old"`
	New int `json:"new"`
}

func (m *natsMailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	request := UpdateMailboxRequest{Old: oldMessages, New: newMessages}
	err = m.conn.StandardRequest("ari.mailboxes.update", name, &request, nil)
	return
}

func (m *natsMailbox) Delete(name string) (err error) {
	err = m.conn.StandardRequest("ari.mailboxes.delete", name, nil, nil)
	return
}
