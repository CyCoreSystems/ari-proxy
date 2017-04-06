package client

/*
type natsSubscription struct {
	m         ari.Matcher
	closeChan chan struct{}
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
	readyCh := make(chan struct{})

	go func() {
		defer sub.Cancel()
		close(readyCh)
		for {
			select {
			case <-ns.closeChan:
				ns.closeChan = nil
				return
			case evt, ok := <-sub.Events():
				if !ok {
					close(ns.closeChan)
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

	<-readyCh
}

func (ns *natsSubscription) Events() chan ari.Event {
	return ns.events
}

func (ns *natsSubscription) Cancel() {
	if ns.closeChan != nil {
		close(ns.closeChan)
	}
}
*/
