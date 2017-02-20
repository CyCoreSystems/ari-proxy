package ariproxy

import "github.com/CyCoreSystems/ari"

type entity struct {
	// Type is the type of entity
	Type string

	// ID is the unique identifier for the entity
	ID string
}

func (s *Server) dialogsForEvent(e ari.Event) (ret []string) {
	for _, i := range entitiesFromEvent(e) {
		ret = append(ret, s.Dialog.List(i.Type, i.ID)...)
	}
	return
}

func entitiesFromEvent(e ari.Event) (ret []entity) {
	if v, ok := e.(ari.ChannelEvent); ok {
		for _, id := range v.GetChannelIDs() {
			ret = append(ret, entity{
				Type: "channel",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(ari.BridgeEvent); ok {
		for _, id := range ev.GetBridgeIDs() {
			ret = append(ret, entity{
				Type: "bridge",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(ari.EndpointEvent); ok {
		for _, id := range ev.GetEndpointIDs() {
			ret = append(ret, entity{
				Type: "endpoint",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(ari.PlaybackEvent); ok {
		for _, id := range ev.GetPlaybackIDs() {
			ret = append(ret, entity{
				Type: "playback",
				ID:   id,
			})
		}
	}

	if ev, ok := e.(ari.RecordingEvent); ok {
		for _, id := range ev.GetRecordingIDs() {
			ret = append(ret, entity{
				Type: "recording",
				ID:   id,
			})
		}
	}

	return
}
