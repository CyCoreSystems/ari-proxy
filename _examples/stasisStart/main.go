package main

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/CyCoreSystems/ari/client/native"

	"github.com/inconshreveable/log15"
)

var ariApp = "test"

var log = log15.New()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	native.Logger = log

	log.Info("Connecting to ARI")
	cl, err := client.New(ctx, client.WithApplication(ariApp))
	if err != nil {
		log.Error("Failed to build ARI client", "error", err)
		return
	}

	// setup app

	log.Info("Starting listener app")

	listenApp(ctx, cl, channelHandler)

	// start call start listener

	log.Info("Starting HTTP Handler")

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make call
		log.Info("Make sample call")
		h, err := createCall(cl)
		if err != nil {
			log.Error("Failed to create call", "error", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Failed to create call: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(h.ID()))
	}))

	log.Info("Listening for requests on port 9990")
	http.ListenAndServe(":9990", nil)

	return
}

func listenApp(ctx context.Context, cl ari.Client, handler func(*ari.ChannelHandle, *ari.StasisStart)) {
	err := client.Listen(ctx, cl, func(ch *ari.ChannelHandle, v *ari.StasisStart) {
		log.Info("Got stasis start", "channel", v.Channel.ID)
		go handler(ch, v)
	})
	if err != nil {
		log.Crit("failed to listen for new calls", "error", err)
	}
	return
}

func createCall(cl ari.Client) (h *ari.ChannelHandle, err error) {
	h, err = cl.Channel().Create(nil, ari.ChannelCreateRequest{
		Endpoint: "Local/1000",
		App:      ariApp,
	})

	return
}

func channelHandler(h *ari.ChannelHandle, startEvent *ari.StasisStart) {
	log.Info("Running channel handler")

	// Subscribe to channel state changes
	stateChange := h.Subscribe(ari.Events.ChannelStateChange)
	defer stateChange.Cancel()

	// Subscribe to StasisEnd events (channel leaving ARI app)
	end := h.Subscribe(ari.Events.StasisEnd)
	defer end.Cancel()

	// Pull the current channel data
	data, err := h.Data()
	if err != nil {
		log.Error("Error getting data", "error", err)
		return
	}
	log.Info("Channel State", "state", data.State)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Info("Waiting for channel events")

		defer wg.Done()

		for {
			select {
			case <-end.Events():
				log.Info("Got stasis end")
				return
			case e := <-stateChange.Events():
				v, ok := e.(*ari.ChannelStateChange)
				if !ok {
					log.Error("failed to interpret event as ChannelStateChange", "error", err)
					return
				}

				log.Info("New Channel State", "state", v.Channel.State)
			}
		}

	}()

	h.Answer()

	wg.Wait()

	h.Hangup()
}
