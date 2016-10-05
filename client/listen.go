package client

import (
	"encoding/json"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Handler defines a function which is called when a new dialog is created
type Handler func(cl *ari.Client, dialog *session.Dialog)

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

			go handler(conn, appStart, h)
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
	// send okay outside of goroutine, so the other side doesn't time out
	data := []byte("ok")
	if err := conn.Publish(reply, data); err != nil {
		Logger.Error("error publishing ok response", "error", err)
		return false
	}

	return true
}

func handler(conn *nats.Conn, appStart session.AppStart, h Handler) {
	d := session.NewDialog(appStart.DialogID, nil)
	d.ChannelID = appStart.ChannelID

	cl, err := New(conn, appStart.Application, d, Options{})
	if err != nil {
		Logger.Error("error creating client", "error", err)
		return
	}

	go func() {
		conn.Subscribe("events.dialog."+d.ID, func(msg *nats.Msg) {
			var ariMessage ari.Message
			ariMessage.SetRaw(&msg.Data)
			Logger.Debug("got eventc", "type", ariMessage.Type)
			cl.Bus.Send(&ariMessage)
		})
	}()

	h(cl, d)
}
