package messagebus

import (
	"strings"
	"time"

	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

// DefaultReconnectionAttemts is the default number of reconnection attempts
// It implements a hard coded fault tolerance for a starting NATS cluster
const DefaultReconnectionAttemts = 5

// DefaultReconnectionWait is the default wating time between each reconnection
// attempt
const DefaultReconnectionWait = 5 * time.Second

// Type is the type of MessageBus (RabbitMQ / NATS)
type Type int

// WildcardType used to identify wildcards used on routing keys on message bus
type WildcardType int

// wildcard types
const (
	WildcardUndefined       WildcardType = iota // undefined type
	WildcardOneWord                             // one word like pre.*.post
	WildcardZeroOrMoreWords                     // zero or more words like pre.>
)

// types
const (
	TypeUnknown  Type = iota // unknown type
	TypeNats                 // NATS type
	TypeRabbitmq             // RabbitMQ type
)

// Server defines the functions used on ari-proxy server
type Server interface {
	Connect() error
	Close()

	SubscribePing(topic string, callback PingHandler) (Subscription, error)
	SubscribeRequest(topic string, callback RequestHandler) (Subscription, error)
	SubscribeRequests(topics []string, callback RequestHandler) (Subscription, error)
	SubscribeCreateRequest(topic string, queue string, callback RequestHandler) (Subscription, error)
	PublishResponse(topic string, msg *proxy.Response) error
	PublishAnnounce(topic string, msg *proxy.Announcement) error
	PublishEvent(topic string, msg ari.Event) error
}

// Client defines the functions used on ari-proxy client
type Client interface {
	Connect() error
	Close()

	SubscribeAnnounce(topic string, callback AnnounceHandler) (Subscription, error)
	SubscribeEvent(topic string, queue string, callback EventHandler) (Subscription, error)

	PublishPing(topic string) error
	Request(topic string, req *proxy.Request) (*proxy.Response, error)
	MultipleRequest(topic string, req *proxy.Request, expectedResp int) ([]*proxy.Response, error)
	MultipleRequestReturnFirstGoodResponse(topic string, req *proxy.Request, expectedResp int) (*proxy.Response, error)

	TimeoutCount() int64
	GetWildcardString(w WildcardType) string
}

// Config has general configuration for MessageBus
type Config struct {
	URL            string
	TimeoutRetries int
	RequestTimeout time.Duration
	ID             string
}

// Subscription defines subscription interface
type Subscription interface {
	Unsubscribe() error
}

// RequestHandler handles requests messages
type RequestHandler func(subject string, reply string, req *proxy.Request)

// ResponseHandler handles response messages
type ResponseHandler func(req *proxy.Response)

// PingHandler handles ping messages
type PingHandler func()

// AnnounceHandler handles announce messages
type AnnounceHandler func(o *proxy.Announcement)

// EventHandler handles event messages
type EventHandler func(b []byte)

// GetType identifies message bus type from an url
func GetType(url string) Type {
	if strings.HasPrefix(url, "amqp://") {
		return TypeRabbitmq
	}
	if strings.HasPrefix(url, "nats://") {
		return TypeNats
	}
	return TypeUnknown
}
