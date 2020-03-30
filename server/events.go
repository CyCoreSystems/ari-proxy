package server

import "github.com/CyCoreSystems/ari/v5"

func (s *Server) dialogsForEvent(e ari.Event) (ret []string) {
	for _, k := range e.Keys() {
		if k == nil {
			s.Log.Debug("received nil key for event", "event", e)
			continue
		}
		ret = append(ret, s.Dialog.List(k.Kind, k.ID)...)
	}
	return
}
