package bus

import (
	"fmt"
	"sync"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"

	"github.com/nats-io/nats.go"
)

// subscriptionEventBufferSize defines the number of events that each
// subscription will queue before accepting more events.
var subscriptionEventBufferSize = 100

// busWrapper binds a NATS subject to an ari.Bus, passing any received NATS messages to that bus
type busWrapper struct {
	subject string

	log log15.Logger

	sub *nats.Subscription

	subs []*subscription // The list of subscriptions

	rwMux sync.RWMutex

	closed bool
}

// A Subscription is a wrapped channel for receiving
// events from the ARI event bus.
type subscription struct {
	key    *ari.Key
	b      *busWrapper     // reference to the event bus
	events []string // list of events to listen for

	mu     sync.Mutex
	closed bool           // channel closure protection flag
	C      chan ari.Event // channel for sending events to the subscriber
}

func newBusWrapper(subject string, nc *nats.EncodedConn, log log15.Logger) (*busWrapper, error) {
	var err error

	w := &busWrapper{
		subject: subject,
		log:     log,
		subs:     []*subscription{},
	}

	w.sub, err = nc.Subscribe(subject, func(m *nats.Msg) {
		w.receive(m)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to subscribe to NATS subject %s", subject)
	}

	return w, nil
}

func (w *busWrapper) receive(o *nats.Msg) {
	var matched bool

	e, err := ari.DecodeEvent(o.Data)
	if err != nil {
		w.log.Error("failed to convert received message to ari.Event", "error", err)
		return
	}

	w.rwMux.RLock()
	// Disseminate the message to the subscribers
	for _, s := range w.subs {
		matched = false
		for _, k := range e.Keys() {
			if matched {
				break
			}

			if s.key.Match(k) {
				matched = true

				for _, topic := range s.events {
					if topic == e.GetType() || topic == ari.Events.All {
						select {
						case s.C <- e:
						default: // never block
						}
					}
				}
			}
		}
	}
	w.rwMux.RUnlock()
}

// Subscribe returns a subscription to the given list
// of event types
func (w *busWrapper) Subscribe(key *ari.Key, eTypes ...string) ari.Subscription {
	s := &subscription{
		key:    key,
		b:      w,
		events: eTypes,
		C:      make(chan ari.Event, subscriptionEventBufferSize),
	}
	w.add(s)
	for k,v := range w.subs {
		fmt.Println(k,v)
	}

	return s
}

// add appends a new subscription to the bus
func (w *busWrapper) add(s *subscription) {
	w.rwMux.Lock()
	w.subs = append(w.subs, s)
	w.rwMux.Unlock()
}

// remove deletes the given subscription from the bus
func (w *busWrapper) remove(s *subscription) {
	w.rwMux.Lock()
	for i, si := range w.subs {
		if s == si {
			// Subs are pointers, so we have to explicitly remove them
			// to prevent memory leaks
			w.subs[i] = w.subs[len(w.subs)-1] // replace the current with the end
			w.subs[len(w.subs)-1] = nil       // remove the end
			w.subs = w.subs[:len(w.subs)-1]   // lop off the end

			break
		}
	}
	w.rwMux.Unlock()
}

func (w *busWrapper) Close() {
	if w.closed {
		return
	}
	if err := w.sub.Unsubscribe(); err != nil {
		w.log.Error("failed to unsubscribe when closing NATS subscription:", err)
	}
	w.closed = true
	for _, s := range w.subs {
		s.Cancel()
	}
}

// Events returns the events channel
func (s *subscription) Events() <-chan ari.Event {
	return s.C
}

// Cancel cancels the subscription and removes it from
// the event bus.
func (s *subscription) Cancel() {
	if s == nil {
		return
	}
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return
	}

	s.closed = true

	s.mu.Unlock()

	// Remove the subscription from the bus
	if s.b != nil {
		s.b.remove(s)
	}

	// Close the subscription's deliver channel
	if s.C != nil {
		close(s.C)
		s.C = nil
	}
	s = nil
}

// Bus provides an ari.Bus interface to NATS
type Bus struct {
	prefix string

	log log15.Logger

	nc *nats.EncodedConn

	subjectBuses map[string]*busWrapper

	mu sync.RWMutex
}

// New returns a new Bus
func New(prefix string, nc *nats.EncodedConn, log log15.Logger) *Bus {
	return &Bus{
		prefix:       prefix,
		log:          log,
		nc:           nc,
		subjectBuses: make(map[string]*busWrapper),
	}
}

type subBus struct {
	bus  ari.Bus
	subs []ari.Subscription
}

func (b *subBus) Close() {
	for _, s := range b.subs {
		s.Cancel()
	}
	b.subs = nil

	// NOTE: we are NOT closing the parent bus here and now, since it could be used by any number of other clients
	// TODO: Ultimately, we will need to derive a way to check to see if the parent bus is then unused, in which case, the NATS subscription(s) should then be closed.
}

func (b *subBus) Send(e ari.Event) {
	b.bus.Send(e)
}

func (b *subBus) Subscribe(key *ari.Key, eTypes ...string) ari.Subscription {
	sub := b.bus.Subscribe(key, eTypes...)

	b.subs = append(b.subs, sub)

	return sub
}

// SubBus creates and returns a new ariBus which is subtended from this one
func (b *Bus) SubBus() ari.Bus {
	return &subBus{
		bus: b,
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

// Close implements ari.Bus
func (b *Bus) Close() {
	b.mu.Lock()

	for _, w := range b.subjectBuses {
		w.Close()
	}

	b.mu.Unlock()
}

// Send implements ari.Bus
func (b *Bus) Send(e ari.Event) {
	// No-op
}

// Subscribe implements ari.Bus
func (b *Bus) Subscribe(key *ari.Key, n ...string) ari.Subscription {
	var err error

	subject := b.subjectFromKey(key)

	b.mu.Lock()
	w, ok := b.subjectBuses[subject]
	if !ok {
		w, err = newBusWrapper(subject, b.nc, b.log)
		if err != nil {
			b.log.Error("failed to create bus wrapper", "key", key, "error", err)
			return nil
		}
		b.subjectBuses[subject] = w
	}
	b.mu.Unlock()

	return w.Subscribe(key, n...)
}
