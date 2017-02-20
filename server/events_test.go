package ariproxy

import (
	"testing"
	"time"

	"github.com/CyCoreSystems/ari"
)

func TestEntitiesFromChannelEvent(t *testing.T) {
	e := ari.ChannelDtmfReceived{
		EventData: ari.EventData{
			Application: "testApp",
			Timestamp:   ari.DateTime(time.Now()),
		},
		Channel: ari.ChannelData{
			ID:   "testChannelEvent",
			Name: "Local/testChannel",
		},
	}

	list := entitiesFromEvent(&e)
	if len(list) != 1 {
		t.Errorf("Incorrect number of entities: %d != 1", len(list))
	}
	if list[0].ID != "testChannelEvent" {
		t.Errorf("Incorrect channel ID: %s != testChannelEvent", list[0].ID)
	}
}
