package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) mailboxData(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Mailbox().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Mailbox: data,
		},
	})
}

func (s *Server) mailboxGet(ctx context.Context, reply string, req *proxy.Request) {
	data, err := s.ari.Mailbox().Data(req.Key)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Key: data.Key,
	})
}

func (s *Server) mailboxDelete(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Mailbox().Delete(req.Key))
}

func (s *Server) mailboxList(ctx context.Context, reply string, req *proxy.Request) {
	list, err := s.ari.Mailbox().List(nil)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.publish(reply, &proxy.Response{
		Keys: list,
	})
}

func (s *Server) mailboxUpdate(ctx context.Context, reply string, req *proxy.Request) {
	s.sendError(reply, s.ari.Mailbox().Update(req.Key, req.MailboxUpdate.Old, req.MailboxUpdate.New))
}
