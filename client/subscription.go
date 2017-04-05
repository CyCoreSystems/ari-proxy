package client

import (
	"sync"

	"github.com/CyCoreSystems/ari"
)

type natsSubscription struct {
	m         ari.Matcher
	closeChan chan struct{}
	events    chan ari.Event
	once      sync.Once
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
	var wg sync.WaitGroup

	readyCh := make(chan struct{})

	wg.Add(1)
	go func() {
		defer sub.Cancel()
		defer ns.Cancel()

		wg.Done()

		for {
			select {
			case <-ns.closeChan:
				return
			case evt, ok := <-sub.Events():
				if !ok {
					return
				}
				if ns.m == nil {
					ns.events <- evt
				} else if ns.m.Match(evt) {
					ns.events <- evt
				}
			}
		}
	}()

	wg.Wait()
}

func (ns *natsSubscription) Events() chan ari.Event {
	return ns.events
}

func (ns *natsSubscription) Cancel() {
	ns.Once(func() {
		close(ns.closeChan)
	})
}
