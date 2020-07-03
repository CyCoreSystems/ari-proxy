package bus

import (
	"fmt"
	"sync"

	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/stdbus"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"

	"github.com/nats-io/nats.go"
)

// busWrapper binds a NATS subject to an ari.Bus, passing any received NATS messages to that bus
type busWrapper struct {
	subject string

	log log15.Logger

	sub *nats.Subscription

	bus ari.Bus
}

func newBusWrapper(subject string, nc *nats.EncodedConn, log log15.Logger) (*busWrapper, error) {
	var err error

	w := &busWrapper{
		subject: subject,
		log:     log,
		bus:     stdbus.New(),
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
	e, err := ari.DecodeEvent(o.Data)
	if err != nil {
		w.log.Error("failed to convert received message to ari.Event", "error", err)
		return
	}

	w.bus.Send(e)
}

func (w *busWrapper) Close() {
	if err := w.sub.Unsubscribe(); err != nil {
		w.log.Error("failed to unsubscribe when closing NATS subscription:", err)
	}
	w.bus.Close()
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
	mu sync.Mutex
}

func (b *subBus) Close() {
	for _, s := range b.subs {
		s.Cancel()
	}
	b.subs = nil

	// NOTE: we are NOT closing the parent bus here and now, since it could be used by any number of other clients
	// TODO: Ultimately, we will need to derive a way to check to see if the parent bus is then unused, in which case, the NATS subscription(s) should then be closed.
}

// Used as callback from stdbus
func (b *subBus) Cancel(s interface{}) {
	b.mu.Lock()
	for i, si := range b.subs {
		if s == si {
			b.subs[i] = b.subs[len(b.subs)-1] // replace the current with the end
			b.subs[len(b.subs)-1] = nil       // remove the end
			b.subs = b.subs[:len(b.subs)-1]   // lop off the end
			break
		}
	}
	b.mu.Unlock()
}

func (b *subBus) Send(e ari.Event) {
	b.bus.Send(e)
}

func (b *subBus) Subscribe(key *ari.Key, eTypes ...string) ari.Subscription {
	sub := b.bus.Subscribe(key, eTypes...)
	sub.AddCancelCallback(b.Cancel)
	b.mu.Lock()
	b.subs = append(b.subs, sub)
	b.mu.Unlock()

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

	return w.bus.Subscribe(key, n...)
}
