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
	ns := newSubscription(p.Get(id))
	ns.Start(p.subscriber, nx...)
	return ns
}
