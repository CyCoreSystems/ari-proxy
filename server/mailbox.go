package ariproxy

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) mailboxData(ctx context.Context, reply string, req *proxy.Request) {
	name := req.MailboxData.Name
	dd, err := s.ari.Mailbox.Data(name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &dd)
}

func (s *Server) mailboxDelete(ctx context.Context, reply string, req *proxy.Request) {
	name := req.MailboxDelete.Name
	err := s.ari.Mailbox.Delete(name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) mailboxList(ctx context.Context, reply string, req *proxy.Request) {
	mx, err := s.ari.Mailbox.List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var mailboxes []string
	for _, m := range mx {
		mailboxes = append(mailboxes, m.ID())
	}

	s.nats.Publish(reply, mailboxes)
}

func (s *Server) mailboxUpdate(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Mailbox.Update(req.MailboxUpdate.Name, req.MailboxUpdate.Old, req.MailboxUpdate.New)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
