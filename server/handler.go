package server

import "github.com/CyCoreSystems/ari-proxy/session"

// Reply is a function which, when called, replies to the request via the
// response object or error.
type Reply func(interface{}, error)

// A Handler2 is a function which provides a session-aware request-response for nats
type Handler2 func(msg *session.Message, reply Reply)

// Handler is left for compat
type Handler func(subj string, request []byte, reply Reply)
