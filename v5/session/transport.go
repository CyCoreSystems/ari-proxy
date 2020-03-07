package session

import "github.com/CyCoreSystems/ari/v5"

// Transport defines how the commands and events are sent.
type Transport interface {

	// Command sends a command and waits for a response
	Command(name string, body interface{}, resp interface{}) error

	// Event dispatches an event
	Event(evt ari.Event) error
}
