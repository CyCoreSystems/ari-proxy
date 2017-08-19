package main

import (
	"golang.org/x/net/context"

	"github.com/inconshreveable/log15"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
	"github.com/CyCoreSystems/ari/ext/record"
)

var ariApp = "test"

var log = log15.New()

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	record.Logger = log

	// connect
	log.Info("Connecting to ARI")
	cl, err := client.New(ctx, client.WithApplication(ariApp), client.WithLogger(log))
	if err != nil {
		log.Error("Failed to build ARI client", "error", err)
		return
	}

	// setup app

	log.Info("Starting listener app")

	err = client.Listen(ctx, cl, appStart(ctx))
	if err != nil {
		log.Error("failed to listen for new calls")
	}
	<-ctx.Done()

	return
}

func appStart(ctx context.Context) func(*ari.ChannelHandle, *ari.StasisStart) {
	return func(h *ari.ChannelHandle, startEvent *ari.StasisStart) {
		defer h.Hangup()

		log.Info("running app", "channel", h.Key().ID)

		if err := h.Answer(); err != nil {
			log.Error("failed to answer call", "error", err)
			//return
		}

		res, err := record.Record(ctx, h,
			record.TerminateOn("any"),
			record.IfExists("overwrite"),
		).Result()
		if err != nil {
			log.Error("failed to record", "error", err)
			return
		}

		if err = res.Save("test-recording"); err != nil {
			log.Error("failed to save recording", "error", err)
		}

		log.Info("completed recording")

		h.Hangup()
	}
}
