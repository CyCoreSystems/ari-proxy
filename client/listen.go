package client

import (
	"context"
	"fmt"

	"github.com/CyCoreSystems/ari-proxy/v5/messagebus"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/rotisserie/eris"
)

// ListenQueue is the queue group to use for distributing StasisStart events to Listeners.
var ListenQueue = "ARIProxyStasisStartDistributorQueue"

// Listen listens for StasisStart events, filtered by the given key.  Any
// matching events will be sent down the returned StasisStart channel.  The
// context which is passed to Listen can be used to stop the Listen execution.
//
// Importantly, the StasisStart events are listened in a NATS/RabbitMQ Queue, which
// means that this may be used to deliver new calls to only a single handler
// out of a set of 1 or more handlers in a cluster.
func Listen(ctx context.Context, ac ari.Client, h func(*ari.ChannelHandle, *ari.StasisStart)) error {
	c, ok := ac.(*Client)
	if !ok {
		return eris.New("ARI Client must be a proxy client")
	}

	subj := fmt.Sprintf(
		"%sevent.%s.%s",
		c.core.prefix,
		c.ApplicationName(),
		c.mbus.GetWildcardString(messagebus.WildcardZeroOrMoreWords),
	)

	c.log.Debug("listening for events", "subject", subj)
	sub, err := c.mbus.SubscribeEvent(subj, ListenQueue, listenProcessor(ac, h))
	if err != nil {
		return eris.Wrap(err, "failed to subscribe to events")
	}
	defer sub.Unsubscribe() // nolint: errcheck

	<-ctx.Done()

	return nil
}

func listenProcessor(ac ari.Client, h func(*ari.ChannelHandle, *ari.StasisStart)) func([]byte) {
	return func(data []byte) {
		e, err := ari.DecodeEvent(data)
		if err != nil {
			Logger.Error("failed to decode event", "error", err)
			return
		}

		Logger.Debug("received event", e.GetType())
		if e.GetType() != "StasisStart" {
			return
		}

		v, ok := e.(*ari.StasisStart)
		if !ok {
			Logger.Error("failed to type-assert StasisStart event")
			return
		}

		h(ari.NewChannelHandle(v.Key(ari.ChannelKey, v.Channel.ID), ac.Channel(), nil), v)
	}
}
