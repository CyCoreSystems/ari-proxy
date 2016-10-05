package ariproxy

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
	uuid "github.com/satori/go.uuid"
)

type creator interface {
	Created() (string, string)
}

type destroyer interface {
	Destroyed() string
}

func (srv *Server) events() {

	if srv.upstream.Bus == nil {
		// useful for tests
		srv.log.Warn("No Upstream Bus in nats event forwarding")
		return
	}

	go func() {
		sub := srv.upstream.Bus.Subscribe(ari.Events.All)
		defer sub.Cancel()

		srv.log.Debug("Listening for events")

		for {
			select {
			case <-srv.ctx.Done():
				return
			case evt := <-sub.Events():

				srv.log.Debug("Got event", "event", evt)

				t := evt.GetType()

				var i *Instance

				switch t {
				case ari.Events.StasisStart:
					i = srv.tryStasisStart(evt)
				case ari.Events.StasisEnd:
					i = srv.tryStasisEnd(evt)
				default:
					i = srv.tryEvent(evt)
				}

				srv.dispatchEvent(evt, i)
			}
		}
	}()
}

func (srv *Server) dispatchEvent(evt ari.Event, i *Instance) {
	raw := *evt.GetRaw()
	app := evt.GetApplication()

	srv.conn.Publish("events.app."+app, raw)

	if i != nil {
		srv.conn.Publish("events.dialog."+i.Dialog.ID, raw)
	}
}

func (srv *Server) tryEvent(evt ari.Event) (i *Instance) {

	c, ok := evt.(creator)
	if ok {
		objectID, relatedID := c.Created()

		i = srv.cache.Find(relatedID)
		if i != nil {
			srv.cache.Add(objectID, i)
		}
	}

	d, ok := evt.(destroyer)
	if ok {
		objectID := d.Destroyed()
		i = srv.cache.Find(objectID)
		if i != nil {
			srv.cache.Remove(objectID, i)
		}
	}

	if i != nil {
		return
	}

	// last resort, check ChannelIDs, BridgeIDs, etc...

	ce, ok := evt.(ari.ChannelEvent)
	if ok {
		for _, ci := range ce.GetChannelIDs() {
			i = srv.cache.Find(ci)
			if i != nil {
				return
			}
		}
	}

	be, ok := evt.(ari.BridgeEvent)
	if ok {
		for _, bi := range be.GetBridgeIDs() {
			i = srv.cache.Find(bi)
			if i != nil {
				return
			}
		}
	}

	ee, ok := evt.(ari.EndpointEvent)
	if ok {
		for _, ei := range ee.GetEndpointIDs() {
			i = srv.cache.Find(ei)
			if i != nil {
				return
			}
		}
	}

	pe, ok := evt.(ari.PlaybackEvent)
	if ok {
		for _, pi := range pe.GetPlaybackIDs() {
			i = srv.cache.Find(pi)
			if i != nil {
				return
			}
		}
	}

	re, ok := evt.(ari.RecordingEvent)
	if ok {
		for _, ri := range re.GetRecordingIDs() {
			i = srv.cache.Find(ri)
			if i != nil {
				return
			}
		}
	}

	return
}

func (srv *Server) tryStasisStart(evt ari.Event) (i *Instance) {
	st := evt.(*ari.StasisStart)

	if i = srv.cache.Find(st.Channel.ID); i != nil {
		return
	}

	// start server side of the component

	id := uuid.NewV1().String()
	i = srv.newInstance(id, nil)
	i.Dialog.ChannelID = st.Channel.ID
	i.Dialog.Objects.Add(st.Channel.ID)
	i.Start(srv.ctx)
	srv.cache.Add(st.Channel.ID, i)

	// send out appstart event

	body, _ := json.Marshal(&session.AppStart{
		ServerID:    srv.ID,
		DialogID:    i.Dialog.ID,
		Application: st.GetApplication(),
		ChannelID:   st.Channel.ID,
	})

	reply := uuid.NewV1().String()

	doneCh := make(chan struct{})
	var err error

	srv.conn.Subscribe(reply, func(msg *nats.Msg) {
		defer close(doneCh)
		body := string(msg.Data)
		if body != "ok" {
			err = errors.New(body)
		}
	})

	srv.conn.PublishRequest("ari.app."+st.GetApplication(), reply, body)

	select {
	case <-doneCh:
	case <-time.After(300 * time.Millisecond):
	}

	//TODO: log error

	return i
}

func (srv *Server) tryStasisEnd(evt ari.Event) (i *Instance) {

	end := evt.(*ari.StasisEnd)

	i = srv.cache.Find(end.Channel.ID)
	if i == nil {
		return
	}

	srv.cache.RemoveAll(i)

	i.Stop()

	return i
}
