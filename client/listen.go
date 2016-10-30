package client

import (
	"context"
	"encoding/json"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
)

// Handler defines a function which is called when a new dialog is created
type Handler func(context.Context, *ari.Client, *session.Dialog)

// Listen listens for an AppStart event and calls the handler when an event comes in
func Listen(ctx context.Context, conn *nats.Conn, appName string, h Handler) error {

	Logger.Debug("Listening on endpoint", "endpoint", "ari.app."+appName)

	ch := make(chan *nats.Msg, 2)
	sub, err := conn.QueueSubscribeSyncWithChan("ari.app."+appName, appName+"_app_listener", ch)
	if err != nil {
		Logger.Debug("Error listening on endpoint", "error", err)
		return errors.Wrap(err, "Unable to subscribe to ARI application start queue")
	}

	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			//TODO: log error
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:

			if !ok {
				return nil
			}
			Logger.Debug("Got message", "msg", string(msg.Data), "ok", ok)

			var appStart session.AppStart
			err := json.Unmarshal(msg.Data, &appStart)
			if err != nil {
				Logger.Error("error unmarshaling appstart", "error", err)
				sendErrorReply(conn, msg.Reply, err)
				continue
			}

			if !sendOkReply(conn, msg.Reply) {
				continue
			}

			go handler(ctx, conn, appStart, h)
		}
	}
}

func sendErrorReply(conn *nats.Conn, reply string, err error) {
	// we got an error in the AppStart, reply with the error
	data := []byte(err.Error())
	if err := conn.Publish(reply, data); err != nil {
		Logger.Error("error publishing error response", "error", err)
	}
}

func sendOkReply(conn *nats.Conn, reply string) bool {
	Logger.Debug("Sending OK", "reply", reply)
	// send okay outside of goroutine, so the other side doesn't time out
	data := []byte("ok")
	if err := conn.Publish(reply, data); err != nil {
		Logger.Error("error publishing ok response", "error", err)
		return false
	}

	return true
}

func handler(ctx context.Context, nc *nats.Conn, appStart session.AppStart, h Handler) {
	// Construct a new Dialog handle
	d := session.NewDialog(appStart.DialogID, nil)
	d.ChannelID = appStart.ChannelID

	// Construct the new ARI client
	cl, err := New(nc, appStart.Application, d, Options{})
	if err != nil {
		Logger.Error("error creating client", "error", err)
		return
	}

	// Bind dialog-related events to the ARI client bus
	sub, err := nc.Subscribe("events.dialog."+d.ID, func(msg *nats.Msg) {
		ariMessage, err := ari.NewMessage(msg.Data)
		if err != nil {
			Logger.Error("failed to create new message from payload", "error", err)
			return
		}

		cl.Bus.Send(ariMessage)
	})
	if err != nil {
		Logger.Error("failed to bind dialog events to ARI client", "error", err)
		return
	}
	defer sub.Unsubscribe()

	// Execute the handler
	h(ctx, cl, d)
}
