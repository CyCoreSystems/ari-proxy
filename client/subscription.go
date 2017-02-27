package client

import "github.com/CyCoreSystems/ari"

type natsSubscription struct {
	m         ari.Matcher
	closeChan chan struct{}
	closed    bool
	events    chan ari.Event
}

func newSubscription(m ari.Matcher) *natsSubscription {
	return &natsSubscription{
		m:         m,
		closeChan: make(chan struct{}),
		events:    make(chan ari.Event, 10),
	}
}

func (ns *natsSubscription) Start(s ari.Subscriber, n ...string) {

	sub := s.Subscribe(n...)

	go func() {
		defer sub.Cancel()
		for {
			select {
			case <-ns.closeChan:
				return
			case evt, ok := <-sub.Events():
				if !ok {
					ns.Cancel()
					continue
				}
				if ns.m == nil {
					ns.events <- evt
				} else if ns.m.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()
}

func (ns *natsSubscription) Events() chan ari.Event {
	return ns.events
}

func (ns *natsSubscription) Cancel() {
	if !ns.closed && ns.closeChan != nil {
		ns.closed = true
		close(ns.closeChan)
	}
}
