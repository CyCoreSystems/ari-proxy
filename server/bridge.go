package ariproxy

import (
	"context"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) bridgeAddChannel(ctx context.Context, reply string, req *proxy.Request) {
	name := req.BridgeAddChannel.Name
	channel := req.BridgeAddChannel.Channel
	err := s.ari.Bridge.AddChannel(name, channel)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, Ok)
}

func (s *Server) bridgeCreate(ctx context.Context, reply string, req *proxy.Request) {

	create := req.BridgeCreate.CreateBridgeRequest
	bh, err := s.ari.Bridge.Create(create.ID, create.Type, create.Name)
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
	bd, err := s.ari.Bridge.Data(req.BridgeData.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &bd)
}

func (s *Server) bridgeDelete(ctx context.Context, reply string, req *proxy.Request) {
	err := s.ari.Bridge.Delete(req.BridgeDelete.Name)
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
	obj, err := s.ari.Bridge.Play(req.BridgePlay.Name, pr.PlaybackID, pr.MediaURI)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &obj)
}

func (s *Server) bridgeRecord(ctx context.Context, reply string, req *proxy.Request) {

	rr := req.BridgeRecord.RecordRequest
	name := req.BridgeRecord.Name

	//TODO: evaluate whether this is needed
	//ins.server.cache.Add(rr.Name, ins)

	var opts ari.RecordingOptions

	opts.Format = rr.Format
	opts.MaxDuration = time.Duration(rr.MaxDuration) * time.Second
	opts.MaxSilence = time.Duration(rr.MaxSilence) * time.Second
	opts.Exists = rr.IfExists
	opts.Beep = rr.Beep
	opts.Terminate = rr.TerminateOn

	obj, err := s.ari.Bridge.Record(name, rr.Name, &opts)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	//NOTE: this originally returned a nil
	s.nats.Publish(reply, &obj)
}

func (s *Server) bridgeRemoveChannel(ctx context.Context, reply string, req *proxy.Request) {
	name := req.BridgeRemoveChannel.Name
	channel := req.BridgeRemoveChannel.Channel

	err := s.ari.Bridge.RemoveChannel(name, channel)
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
