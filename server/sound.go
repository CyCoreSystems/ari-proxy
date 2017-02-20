package ariproxy

import (
	"context"

	"github.com/CyCoreSystems/ari-proxy/proxy"
)

func (s *Server) soundData(ctx context.Context, reply string, req *proxy.Request) {
	sd, err := s.ari.Sound.Data(req.SoundData.Name)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	s.nats.Publish(reply, &sd)
}

func (s *Server) soundList(ctx context.Context, reply string, req *proxy.Request) {

	filters := req.SoundList.Filters

	if len(filters) == 0 {
		filters = nil // just send nil to upstream if empty. makes tests easier
	}

	sx, err := s.ari.Sound.List(filters)
	if err != nil {
		s.sendError(reply, err)
		return
	}

	var sounds []string
	for _, sound := range sx {
		sounds = append(sounds, sound.ID())
	}

	s.nats.Publish(reply, sounds)
}
