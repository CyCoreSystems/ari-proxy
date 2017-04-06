package client

import "github.com/nats-io/nats"

// Conn is the wrapper type for a nats connnection along some ARI specific options
type Conn struct {
}

// ReadRequest sends a request that is a "read"... a request
// which can be retried as needed without consequence.
// NOTE: It is less about "read" operations and more about
// operations which are repeatable/idempotent.
func (c *Conn) ReadRequest(cmd string, name string, body interface{}, dest interface{}) (err error) {
	panic("Removed")
}

// StandardRequest is a request that sends JSON and receives JSON (on success)
// OR receives an error from the remote server
func (c *Conn) StandardRequest(cmd string, name string, body interface{}, dest interface{}) (err error) {
	panic("Removed")
}

// RawRequest sends a tiered request, with an initial OK to acknowledge receipt,
// followed by either a response or an error.
func (c *Conn) RawRequest(cmd string, name string, data []byte) (msg *nats.Msg, err error) {
	panic("Removed")
}

// --

type temp interface {
	Temporary() bool
}

type timeoutErr string

func (err timeoutErr) Error() string {
	return string(err)
}

func (err timeoutErr) Timeout() bool {
	return true
}

func (err timeoutErr) Temporary() bool {
	return true
}
