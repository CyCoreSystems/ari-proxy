package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type mailbox struct {
	c *Client
}

func (m *mailbox) Get(key *ari.Key) *ari.MailboxHandle {
	k, err := m.c.getRequest(&proxy.Request{
		Kind: "MailboxGet",
		Key:  key,
	})
	if err != nil {
		m.c.log.Warn("failed to get bridge for handle", "error", err)
		return ari.NewMailboxHandle(key, m)
	}
	return ari.NewMailboxHandle(k, m)
}

func (m *mailbox) List(filter *ari.Key) ([]*ari.Key, error) {
	return m.c.listRequest(&proxy.Request{
		Kind: "MailboxList",
		Key:  filter,
	})
}

func (m *mailbox) Data(key *ari.Key) (*ari.MailboxData, error) {
	data, err := m.c.dataRequest(&proxy.Request{
		Kind: "MailboxData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.Mailbox, nil
}

func (m *mailbox) Update(key *ari.Key, oldMessages int, newMessages int) error {
	return m.c.commandRequest(&proxy.Request{
		Kind: "MailboxUpdate",
		Key:  key,
		MailboxUpdate: &proxy.MailboxUpdate{
			New: newMessages,
			Old: oldMessages,
		},
	})
}

func (m *mailbox) Delete(key *ari.Key) error {
	return m.c.commandRequest(&proxy.Request{
		Kind: "MailboxDelete",
		Key:  key,
	})
}
