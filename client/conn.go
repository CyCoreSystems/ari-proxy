package client

import (
	"encoding/json"
	"time"

	"github.com/CyCoreSystems/ari-proxy/session"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Conn is the wrapper type for a nats connnection along some ARI specific options
type Conn struct {
	opts        Options
	conn        *nats.Conn
	dialog      *session.Dialog
	application string
}

// ReadRequest sends a request that is a "read"... a request
// which can be retried as needed without consequence.
// NOTE: It is less about "read" operations and more about
// operations which are repeatable/idempotent.
func (c *Conn) ReadRequest(cmd string, name string, body interface{}, dest interface{}) (err error) {

	maxRetries := c.opts.ReadOperationRetryCount

	if maxRetries == 0 {
		maxRetries = 1
	}

	for i := 0; i <= maxRetries; i++ {

		err = c.StandardRequest(cmd, name, body, dest)
		if err == nil {
			return
		}

		if t, ok := err.(temp); !ok || !t.Temporary() {
			return
		}
	}

	return
}

// StandardRequest is a request that sends JSON and receives JSON (on success)
// OR receives an error from the remote server
func (c *Conn) StandardRequest(cmd string, name string, body interface{}, dest interface{}) (err error) {

	// build json request

	data := []byte("{}")

	if data != nil {
		if data, err = json.Marshal(body); err != nil {
			return
		}
	}

	// make request

	var msg *nats.Msg
	msg, err = c.RawRequest(cmd, name, data)

	// parse json response

	if msg != nil && dest != nil {
		err = json.Unmarshal(msg.Data, dest)
	}

	return
}

// RawRequest sends a tiered request, with an initial OK to acknowledge receipt,
// followed by either a response or an error.
func (c *Conn) RawRequest(cmd string, name string, data []byte) (msg *nats.Msg, err error) {

	conn := c.conn

	// build message

	var subj string
	if c.dialog != nil {
		subj = "ari.commands.dialog." + c.dialog.ID
	} else {
		subj = "ari.commands.dialog." + c.application + "_0"
	}

	var sessionMessage session.Message
	sessionMessage.Command = cmd
	sessionMessage.Object = name
	sessionMessage.Payload = data

	// convert to json

	var body []byte
	body, err = json.Marshal(&sessionMessage)
	if err != nil {
		return
	}

	// prepare response channel

	var sub *nats.Subscription

	replyID := uuid.NewV1().String()
	ch := make(chan *nats.Msg, 2)
	sub, err = conn.ChanSubscribe(replyID+".>", ch)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	// send request

	if err = conn.PublishRequest(subj, replyID, body); err != nil {
		return
	}

	requestTimeout := c.opts.RequestTimeout

	respType := "ok"

	for {
		// listen for response

		select {
		case msg = <-ch:
		case <-time.After(requestTimeout):
			err = timeoutErr("Timeout waiting for response type " + respType)
			return
		}

		// handle err or "OK, keep waiting" from server
		msgType := msg.Subject[len(replyID)+1:]

		switch msgType {
		case "err":
			data := msg.Data
			msg = nil // zero out msg on error
			m := make(map[string]interface{})
			if err2 := json.Unmarshal(data, &m); err2 != nil {
				err = errors.Wrap(err2, "Error decoding remote error")
			} else {
				err = &remoteError{subj, MapToError(m)}
			}
			return
		case "ok":
			requestTimeout = 2 * time.Second
			respType = "body"
			continue
		default:
			return
		}
	}
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
