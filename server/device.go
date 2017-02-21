package server

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) deviceStateData(ctx context.Context, reply string, req *proxy.Request) {
	id := req.DeviceStateData.ID
	dd, err := s.ari.DeviceState.Data(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}
	s.nats.Publish(reply, &dd)
}

func (s *Server) deviceStateDelete(ctx context.Context, reply string, req *proxy.Request) {
	id := req.DeviceStateDelete.ID
	err := s.ari.DeviceState.Delete(id)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}

func (s *Server) deviceStateList(ctx context.Context, reply string, req *proxy.Request) {
	dx, err := s.ari.DeviceState.List()
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var ret []string
	for _, device := range dx {
		ret = append(ret, device.ID())
	}

	s.nats.Publish(reply, ret)
}

func (s *Server) deviceStateUpdate(ctx context.Context, reply string, req *proxy.Request) {
	id := req.DeviceStateUpdate.ID
	state := req.DeviceStateUpdate.State
	err := s.ari.DeviceState.Update(id, state)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.sendError(reply, nil)
}
