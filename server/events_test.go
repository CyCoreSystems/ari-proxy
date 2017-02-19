package ariproxy

import (
	"fmt"
	"testing"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
)

// tests for incoming event handling

func TestTryEventNoAssociations(t *testing.T) {
	var srv Server

	ix := srv.tryEvent(&ari.PlaybackStarted{
		Playback: ari.PlaybackData{
			ID:        "pb1",
			TargetURI: "channel:channel1",
		},
	})

	if len(ix) != 0 {
		t.Errorf("List for non-associationed playback event should be 0")
	}
}

func TestTryEventOneAssociation(t *testing.T) {
	var srv Server
	srv.cache.Init()

	var i Instance
	i.Dialog = session.NewDialog("d1", nil)
	srv.cache.Add("channel1", &i)

	ix := srv.tryEvent(&ari.PlaybackStarted{
		Playback: ari.PlaybackData{
			ID:        "pb1",
			TargetURI: "channel:channel1",
		},
	})

	if len(ix) != 1 {
		t.Errorf("List for playback event should be 1, was '%d'", len(ix))
	}

	fmt.Printf("%v\n", ix)
}

func TestTryEventMultipleAssociations(t *testing.T) {
	var srv Server
	srv.cache.Init()

	var i1 Instance
	i1.Dialog = session.NewDialog("d1", nil)
	var i2 Instance
	i2.Dialog = session.NewDialog("d2", nil)

	srv.cache.Add("channel1", &i1)
	srv.cache.Add("bridge1", &i1)

	srv.cache.Add("channel2", &i2)

	ix := srv.tryEvent(&ari.ChannelEnteredBridge{
		Channel: ari.ChannelData{
			ID: "channel2",
		},
		Bridge: ari.BridgeData{
			ID: "bridge1",
		},
	})

	if len(ix) != 2 {
		t.Errorf("List for channel entered bridge should be 2, was '%d'", len(ix))
	}

	fmt.Printf("%v\n", ix)
}
