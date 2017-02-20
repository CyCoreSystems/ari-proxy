package ariproxy

import (
	"context"
	"time"

	"github.com/CyCoreSystems/ari"
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

	err := s.ari.Bridge.AddChannel(id, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, Ok)
}

func (s *Server) bridgeCreate(ctx context.Context, reply string, req *proxy.Request) {

	create := req.BridgeCreate.CreateBridgeRequest

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", create.ID)
	}

	bh, err := s.ari.Bridge.Create(create.ID, create.Type, create.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//TODO: evaluate how to perform this
	// and if we need it
	//s.cache.Add(create.ID, ins)

	s.nats.Publish(reply, bh.ID())
}

func (s *Server) bridgeData(ctx context.Context, reply string, req *proxy.Request) {
	bd, err := s.ari.Bridge.Data(req.BridgeData.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &bd)
}

func (s *Server) bridgeDelete(ctx context.Context, reply string, req *proxy.Request) {
	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgeDelete.ID)
	}

	err := s.ari.Bridge.Delete(req.BridgeDelete.ID)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, Ok)
}

func (s *Server) bridgeList(ctx context.Context, reply string, req *proxy.Request) {
	bx, err := s.ari.Bridge.List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var bridges []string
	for _, bridge := range bx {
		bridges = append(bridges, bridge.ID())
	}

	s.nats.Publish(reply, &bridges)
}

func (s *Server) bridgePlay(ctx context.Context, reply string, req *proxy.Request) {

	pr := req.BridgePlay.PlayRequest

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", req.BridgePlay.ID)
		s.Dialog.Bind(req.Metadata.Dialog, "playback", pr.PlaybackID)
	}

	obj, err := s.ari.Bridge.Play(req.BridgePlay.ID, pr.PlaybackID, pr.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &obj)
}

func (s *Server) bridgeRecord(ctx context.Context, reply string, req *proxy.Request) {

	rr := req.BridgeRecord.RecordRequest
	id := req.BridgeRecord.ID

	//TODO: evaluate whether this is needed
	//ins.server.cache.Add(rr.Name, ins)

	var opts ari.RecordingOptions

	opts.Format = rr.Format
	opts.MaxDuration = time.Duration(rr.MaxDuration) * time.Second
	opts.MaxSilence = time.Duration(rr.MaxSilence) * time.Second
	opts.Exists = rr.IfExists
	opts.Beep = rr.Beep
	opts.Terminate = rr.TerminateOn

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", id)
		s.Dialog.Bind(req.Metadata.Dialog, "recording", rr.Name)
	}

	obj, err := s.ari.Bridge.Record(id, rr.Name, &opts)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &obj)
}

func (s *Server) bridgeRemoveChannel(ctx context.Context, reply string, req *proxy.Request) {
	id := req.BridgeRemoveChannel.ID
	channel := req.BridgeRemoveChannel.Channel

	// bind dialog
	if req.Metadata.Dialog != "" {
		s.Dialog.Bind(req.Metadata.Dialog, "bridge", id)
		s.Dialog.Bind(req.Metadata.Dialog, "channel", channel) //TODO: do we unbind here? probably not
	}

	err := s.ari.Bridge.RemoveChannel(id, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, Ok)
}

func (s *Server) bridgeSubscribe(ctx context.Context, reply string, req *proxy.Request) {

	// check for existence
	_, err := s.ari.Bridge.Data(req.BridgeSubscribe.ID)
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
