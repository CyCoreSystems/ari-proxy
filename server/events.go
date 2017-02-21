package server

import "github.com/CyCoreSystems/ari"

func (s *Server) dialogsForEvent(e ari.Event) (ret []string) {
	for _, i := range ari.EntitiesFromEvent(e) {
		ret = append(ret, s.Dialog.List(i.Type, i.ID)...)
	}
	return
}
