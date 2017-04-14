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
	req := proxy.Request{
		MailboxList: &proxy.MailboxList{},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.GetSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		mx = append(mx, m.Get(i.ID))
	}
	return
}

func (m *mailbox) Data(name string) (d *ari.MailboxData, err error) {
	req := proxy.Request{
		MailboxData: &proxy.MailboxData{
			Name: name,
		},
	}
	var resp proxy.DataResponse
	err = m.c.nc.Request(proxy.GetSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	d = resp.MailboxData
	return
}

func (m *mailbox) Update(name string, oldMessages int, newMessages int) (err error) {
	req := proxy.Request{
		MailboxUpdate: &proxy.MailboxUpdate{
			Name: name,
			Old:  oldMessages,
			New:  newMessages,
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

func (m *mailbox) Delete(name string) (err error) {
	req := proxy.Request{
		MailboxDelete: &proxy.MailboxDelete{
			Name: name,
		},
	}
	var resp proxy.Response
	err = m.c.nc.Request(proxy.CommandSubject(m.c.prefix, m.c.appName, ""), &req, &resp, m.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
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
