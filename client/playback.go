package client

import "github.com/CyCoreSystems/ari"

type natsPlayback struct {
	conn       *Conn
	subscriber ari.Subscriber
}

func (p *natsPlayback) Get(id string) *ari.PlaybackHandle {
	return ari.NewPlaybackHandle(id, p)
}

func (p *natsPlayback) Data(id string) (d ari.PlaybackData, err error) {
	err = p.conn.ReadRequest("ari.playback.data", id, nil, &d)
	return
}

func (p *natsPlayback) Control(id string, op string) (err error) {
	err = p.conn.StandardRequest("ari.playback.control", id, &op, nil)
	return
}

func (p *natsPlayback) Stop(id string) (err error) {
	err = p.conn.StandardRequest("ari.playback.stop", id, nil, nil)
	return
}

func (p *natsPlayback) Subscribe(id string, nx ...string) ari.Subscription {

	var ns natsSubscription

	ns.events = make(chan ari.Event, 10)
	ns.closeChan = make(chan struct{})

	playbackHandle := p.Get(id)

	go func() {
		sub := p.subscriber.Subscribe(nx...)
		defer sub.Cancel()
		for {

			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt := <-sub.Events():
				if playbackHandle.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()

	return &ns
}
