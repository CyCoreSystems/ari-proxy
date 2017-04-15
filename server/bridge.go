package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) bridgeAddChannel(ctx context.Context, reply string, req *proxy.Request) {
	id := req.BridgeAddChannel.ID
	channel := req.BridgeAddChannel.Channel

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", id)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", channel)
	}

	err := s.ari.Bridge().AddChannel(id, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeCreate(ctx context.Context, reply string, req *proxy.Request) {

	create := req.BridgeCreate

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", create.ID)
	}

	bh, err := s.ari.Bridge().Create(req.BridgeCreate.ID, req.BridgeCreate.Type, req.BridgeCreate.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Entity: &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       bh.ID(),
		},
	})
}

func (s *Server) bridgeData(ctx context.Context, reply string, req *proxy.Request) {
	bd, err := s.ari.Bridge().Data(req.BridgeData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &proxy.Response{
		Data: &proxy.EntityData{
			Metadata: s.Metadata(req.Metadata.Dialog),
			Bridge:   bd,
		},
	})
}

func (s *Server) bridgeDelete(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgeDelete.ID)
	}

	err := s.ari.Bridge().Delete(req.BridgeDelete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeList(ctx context.Context, reply string, req *proxy.Request) {
	bx, err := s.ari.Bridge().List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	resp := proxy.EntityList{List: []*proxy.Entity{}}
	for _, i := range bx {
		resp.List = append(resp.List, &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       i.ID(),
		})
	}

	s.nats.Publish(reply, &proxy.Response{
		EntityList: &resp,
	})
}

func (s *Server) bridgePlay(ctx context.Context, reply string, req *proxy.Request) {

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgePlay.ID)
		s.Dialog.Bind(req.Metadata.Dialog, "playback", req.BridgePlay.PlaybackID)
	}

	obj, err := s.ari.Bridge().Play(req.BridgePlay.ID, req.BridgePlay.PlaybackID, req.BridgePlay.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &proxy.Response{
		Entity: &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       obj.ID(),
		},
	})
}

func (s *Server) bridgeRecord(ctx context.Context, reply string, req *proxy.Request) {

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgeRecord.ID)
		s.Dialog.Bind(req.Metadata.Dialog, "recording", req.BridgeRecord.Name)
	}

	obj, err := s.ari.Bridge().Record(req.BridgeRecord.ID, req.BridgeRecord.Name, req.BridgeRecord.Options)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &proxy.Response{
		Entity: &proxy.Entity{
			Metadata: s.Metadata(req.Metadata.Dialog),
			ID:       obj.ID(),
		},
	})
}

func (s *Server) bridgeRemoveChannel(ctx context.Context, reply string, req *proxy.Request) {
	id := req.BridgeRemoveChannel.ID
	channel := req.BridgeRemoveChannel.Channel

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", id)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", channel) //TODO: do we unbind here? probably not
	}

	err := s.ari.Bridge().RemoveChannel(id, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) bridgeSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// check for existence
	_, err := s.ari.Bridge().Data(req.BridgeSubscribe.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgeSubscribe.ID)
	}

	s.sendError(reply, nil)
}
