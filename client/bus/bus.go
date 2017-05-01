package bus

import (
	"fmt"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari"
	"github.com/nats-io/nats"
)

// Bus provides an ari.Bus interface to NATS
type Bus struct {
	prefix string

	log log15.Logger

	nc *nats.EncodedConn
}

// New returns a new Bus
func New(prefix string, nc *nats.EncodedConn, log log15.Logger) *Bus {
	return &Bus{
		prefix: prefix,
		log:    log,
		nc:     nc,
	}
}

func (b *Bus) subjectFromKey(key *ari.Key) string {
	if key == nil {
		return fmt.Sprintf("%sevent.>", b.prefix)
	}

	if key.Dialog != "" {
		return fmt.Sprintf("%sdialogevent.%s", b.prefix, key.Dialog)
	}

	subj := fmt.Sprintf("%sevent.", b.prefix)
	if key.App == "" {
		return subj + ">"
	}
	subj += key.App + "."

	if key.Node == "" {
		return subj + ">"
	}
	return subj + key.Node
}

// Subscription represents an ari.Subscription over NATS
type Subscription struct {
	key *ari.Key

	log log15.Logger

	subscription *nats.Subscription

	eventChan chan ari.Event

	events []string

	closed bool
}

// Close implements ari.Bus
func (b *Bus) Close() {
	return
}

// Send implements ari.Bus
func (b *Bus) Send(e ari.Event) {
	return
}

// Subscribe implements ari.Bus
func (b *Bus) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	var err error

	s := &Subscription{
		key:       key,
		log:       b.log,
		eventChan: make(chan ari.Event),
		events:    n,
	}

	s.subscription, err = b.nc.Subscribe(b.subjectFromKey(key), func(m *nats.Msg) {
		s.receive(m)
	})
	if err != nil {
		b.log.Error("failed to subscribe to NATS", "error", err)
		return nil
	}
	return s
}

// Events returns the channel on which events from this subscription will be sent
func (s *Subscription) Events() <-chan ari.Event {
	return s.eventChan
}

// Cancel destroys the subscription
func (s *Subscription) Cancel() {
	if s.subscription != nil {
		err := s.subscription.Unsubscribe()
		if err != nil {
			s.log.Error("failed unsubscribe from NATS", "error", err)
		}
	}

	if !s.closed {
		s.closed = true
		close(s.eventChan)
	}
}

func (s *Subscription) receive(o *nats.Msg) {
	e, err := ari.DecodeEvent(o.Data)
	if err != nil {
		s.log.Error("failed to convert received message to ari.Event", "error", err)
		return
	}

	if s.matchEvent(e) {
		if !s.closed {
			s.eventChan <- e
		}
	}
}

func (s *Subscription) matchEvent(o ari.Event) bool {
	// First, filter by type
	var match bool
	for _, kind := range s.events {
		if kind == o.GetType() || kind == ari.Events.All {
			match = true
			break
		}
	}
	if !match {
		return false
	}

	// If we don't have a resource ID, we match everything
	// Next, match the entity
	if s.key == nil || s.key.ID != "" {
		return true
	}

	for _, k := range o.Keys() {
		if s.key.Match(k) {
			return true
		}
	}
	return false
}
