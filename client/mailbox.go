package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type mailbox struct {
	c *Client
}

func (m *mailbox) Get(name string) ari.MailboxHandle {
	return &mailboxHandle{
		id:      name,
		mailbox: m,
	}
}

func (m *mailbox) List() (mx []ari.MailboxHandle, err error) {
	ml, err := m.c.listRequest(&proxy.Request{
		MailboxList: &proxy.MailboxList{},
	})
	if err != nil {
		return
	}
	for _, i := range ml.List {
		mx = append(mx, m.Get(i.ID))
	}
	return
}

func (m *mailbox) Data(name string) (d *ari.MailboxData, err error) {
	data, err := m.c.dataRequest(&proxy.Request{
		MailboxData: &proxy.MailboxData{
			Name: name,
		},
	})
	if err != nil {
		return
	}
	d = data.Mailbox
	return
}

func (m *mailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	err = m.c.commandRequest(&proxy.Request{
		MailboxUpdate: &proxy.MailboxUpdate{
			Name: name,
			Old:  oldMessages,
			New:  newMessages,
		},
	})
	return
}

func (m *mailbox) Delete(name string) (err error) {
	err = m.c.commandRequest(&proxy.Request{
		MailboxDelete: &proxy.MailboxDelete{
			Name: name,
		},
	})
	return
}

type mailboxHandle struct {
	id      string
	mailbox *mailbox
}

func (m *mailboxHandle) Data() (*ari.MailboxData, error) {
	return m.mailbox.Data(m.id)
}

func (m *mailboxHandle) Delete() error {
	return m.mailbox.Delete(m.id)
}

func (m *mailboxHandle) ID() string {
	return m.id
}

func (m *mailboxHandle) Update(old int, new int) error {
	return m.mailbox.Update(m.id, old, new)
}
