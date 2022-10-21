package bus

import (
	"fmt"
	"sync"

	"github.com/CyCoreSystems/ari-proxy/v5/messagebus"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/inconshreveable/log15"
)

// EventChanBufferLength is the number of unhandled events which can be queued
// to the event channel buffer before further events are lost.
var EventChanBufferLength = 10

// Bus provides an ari.Bus interface to MessageBus
type Bus struct {
	prefix string

	log log15.Logger

	mbus messagebus.Client
}

// New returns a new Bus
func New(prefix string, m messagebus.Client, log log15.Logger) *Bus {
	return &Bus{
		prefix: prefix,
		log:    log,
		mbus:   m,
	}
}

func (b *Bus) subjectFromKey(key *ari.Key) string {
	if key == nil {
		return fmt.Sprintf(
			"%sevent.%s",
			b.prefix,
			b.mbus.GetWildcardString(messagebus.WildcardZeroOrMoreWords),
		)
	}

	if key.Dialog != "" {
		return fmt.Sprintf("%sdialogevent.%s", b.prefix, key.Dialog)
	}

	subj := fmt.Sprintf("%sevent.", b.prefix)
	if key.App == "" {
		return subj + b.mbus.GetWildcardString(messagebus.WildcardZeroOrMoreWords)
	}
	subj += key.App + "."

	if key.Node == "" {
		return subj + b.mbus.GetWildcardString(messagebus.WildcardZeroOrMoreWords)
	}
	return subj + key.Node
}

// Subscription represents an ari.Subscription over MessageBus
type Subscription struct {
	key *ari.Key

	log log15.Logger

	subscription messagebus.Subscription

	eventChan chan ari.Event

	events []string

	closed bool

	mu sync.RWMutex
}

// Close implements ari.Bus
func (b *Bus) Close() {
	// No-op
}

// Send implements ari.Bus
func (b *Bus) Send(e ari.Event) {
	// No-op
}

// Subscribe implements ari.Bus
func (b *Bus) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	var err error

	s := &Subscription{
		key:       key,
		log:       b.log,
		eventChan: make(chan ari.Event, EventChanBufferLength),
		events:    n,
	}

	var app string
	if key != nil {
		app = key.App
	}

	s.subscription, err = b.mbus.SubscribeEvent(
		b.subjectFromKey(key),
		app,
		s.receive,
	)
	if err != nil {
		b.log.Error("failed to subscribe to MessageBus", "error", err)
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
	if s == nil {
		return
	}

	if s.subscription != nil {
		err := s.subscription.Unsubscribe()
		if err != nil {
			s.log.Error("failed unsubscribe from MessageBus", "error", err)
		}
	}

	s.mu.Lock()
	if !s.closed {
		s.closed = true
		close(s.eventChan)
	}
	s.mu.Unlock()
}

func (s *Subscription) receive(data []byte) {
	e, err := ari.DecodeEvent(data)
	if err != nil {
		s.log.Error("failed to convert received message to ari.Event", "error", err)
		return
	}

	if s.matchEvent(e) {
		s.mu.RLock()
		if !s.closed {
			s.eventChan <- e
		}
		s.mu.RUnlock()
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
	if s.key == nil || s.key.ID == "" {
		return true
	}

	for _, k := range o.Keys() {
		if s.key.Match(k) {
			return true
		}
	}
	return false
}
