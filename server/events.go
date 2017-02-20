package ariproxy

/*
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

				var ix []*Instance

				switch t {
				case ari.Events.StasisStart:
					ix = srv.tryStasisStart(evt)
				case ari.Events.StasisEnd:
					ix = srv.tryStasisEnd(evt)
				default:
					ix = srv.tryEvent(evt)
				}

				for _, i := range ix {
					srv.dispatchEvent(evt, i)
				}
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

func (srv *Server) tryEvent(evt ari.Event) []*Instance {
	var m = make(map[string]*Instance)

	c, ok := evt.(creator)
	if ok {
		objectID, relatedID := c.Created()

		for _, i := range srv.cache.Find(relatedID) {
			srv.cache.Add(objectID, i)
			m[i.Dialog.ID] = i
		}

		// if the created event has a list of channel IDs and
		// those channels have instances, associate the the created object
		ce, ok := evt.(ari.ChannelEvent)
		if ok {
			for _, ci := range ce.GetChannelIDs() {
				for _, i := range srv.cache.Find(ci) {
					srv.cache.Add(objectID, i)
					m[i.Dialog.ID] = i
				}
			}
		}
	}

	// last resort, check ChannelIDs, BridgeIDs, etc...

	ce, ok := evt.(ari.ChannelEvent)
	if ok {
		for _, ci := range ce.GetChannelIDs() {
			for _, i := range srv.cache.Find(ci) {
				m[i.Dialog.ID] = i
			}
		}
	}

	be, ok := evt.(ari.BridgeEvent)
	if ok {
		for _, bi := range be.GetBridgeIDs() {
			for _, i := range srv.cache.Find(bi) {
				m[i.Dialog.ID] = i
			}
		}
	}

	ee, ok := evt.(ari.EndpointEvent)
	if ok {
		for _, ei := range ee.GetEndpointIDs() {
			for _, i := range srv.cache.Find(ei) {
				m[i.Dialog.ID] = i
			}
		}
	}

	pe, ok := evt.(ari.PlaybackEvent)
	if ok {
		for _, pi := range pe.GetPlaybackIDs() {
			for _, i := range srv.cache.Find(pi) {
				m[i.Dialog.ID] = i
			}
		}
	}

	re, ok := evt.(ari.RecordingEvent)
	if ok {
		for _, ri := range re.GetRecordingIDs() {
			for _, i := range srv.cache.Find(ri) {
				m[i.Dialog.ID] = i
			}
		}
	}

	var ix []*Instance
	for _, i := range m {
		ix = append(ix, i)
	}

	return ix
}

func (srv *Server) tryStasisStart(evt ari.Event) (il []*Instance) {
	st := evt.(*ari.StasisStart)

	// start server side of the component

	srv.log.Debug("Sending out AppStart to endpoint", "endpoint", "ari.app."+st.GetApplication())

	id := uuid.NewV1().String()
	i := srv.newInstance(id, nil)
	i.Dialog.ChannelID = st.Channel.ID
	i.Dialog.Objects.Add(st.Channel.ID)
	i.Start(srv.ctx)
	srv.cache.Add(st.Channel.ID, i)

	// send out appstart event

	body, _ := json.Marshal(&session.AppStart{
		ServerID:    srv.ID,
		DialogID:    i.Dialog.ID,
		Application: st.GetApplication(),
		AppArgs:     st.Args,
		ChannelID:   st.Channel.ID,
	})

	reply := uuid.NewV1().String()

	var err error

	ch := make(chan *nats.Msg, 1)
	defer close(ch)
	sub, err := srv.conn.ChanSubscribe(reply, ch)
	if err != nil {
		srv.log.Error("Error subscribing on reply channel", "error", err)
		return
	}
	defer sub.Unsubscribe()

	srv.log.Debug("Publishing stasis start request", "appname", st.GetApplication())

	if err = srv.conn.PublishRequest("ari.app."+st.GetApplication(), reply, body); err != nil {
		srv.log.Error("Error publishing StasisStart request", "error", err, "appname", st.GetApplication())
		return
	}

	il = append(il, i)

	srv.log.Debug("Waiting for response")

	select {
	case msg := <-ch:
		body := string(msg.Data)
		if body != "ok" {
			err = errors.New(body)
		}
		srv.log.Debug("Got response", "error", err)
	case <-time.After(1000 * time.Millisecond):
		srv.log.Error("Timed out")
	}

	if err != nil {
		srv.log.Error("Error in StasisStart Handler", "error", err)
	}

	return
}

func (srv *Server) tryStasisEnd(evt ari.Event) (il []*Instance) {

	end := evt.(*ari.StasisEnd)

	il = srv.cache.Find(end.Channel.ID)

	for _, i := range il {
		srv.cache.RemoveObject(end.Channel.ID, i)
		if len(i.Dialog.Objects.Items()) == 0 {
			i.Stop()
		}
	}

	return
}
*/
