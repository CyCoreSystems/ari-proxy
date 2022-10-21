# ari-proxy
[![Build Status](https://travis-ci.org/CyCoreSystems/ari-proxy.png)](https://travis-ci.org/CyCoreSystems/ari-proxy) [![](https://godoc.org/github.com/CyCoreSystems/ari-proxy?status.svg)](https://godoc.org/github.com/CyCoreSystems/ari-proxy)

Proxy for the Asterisk REST interface (ARI).

The ARI proxy facilitates scaling of both applications and Asterisk,
independently and with minimal coordination.  Each Asterisk instance and ARI
application pair runs an `ari-proxy` server instance, which talks to a common
NATS or RabbitMQ cluster.  Each client application talks to the same message bus.  The
clients automatically and continuously discover new Asterisk instances, so the
only coordination needed is the common location of the message bus.

The ARI proxy allows for:
  - Any number of applications running the ARI client
  - Any number of `ari-proxy` services running on any number of Asterisk
    instances
  - Simple call control throughout the cluster, regardless of which Asterisk
    instance is controlling the call
  - Simple call distribution regardless of the number of potential application
    services.  (New calls are automatically sent to a single recipient
    application.)
  - Simple call event reception by any number of application clients.  (No
    single-app lockout)

Supported message buses:
  - [NATS](https://nats.io)
  - [RabbitMQ](https://rabbitmq.com)

## Proxy server


Docker images are kept up to date with releases and are tagged accordingly.  The
`ari-proxy` does not expose any services, so no ports need to be opened for it.
However, it does need to know how to connect to both Asterisk and the message
bus.

```
   docker run \
     -e ARI_APPLICATION="my_test_app" \
     -e ARI_USERNAME="demo-user" \
     -e ARI_PASSWORD="supersecret" \
     -e ARI_HTTP_URL="http://asterisk:8088/ari" \
     -e ARI_WEBSOCKET_URL="ws://asterisk:8088/ari/events" \
     -e MESSAGEBUS_URL="nats://nats:4222" \
     cycoresystems/ari-proxy
```

Binary releases are available on the [releases page](https://github.com/CyCoreSystems/ari-proxy/releases).

You can also install the server manually:

```
   go install github.com/CyCoreSystems/ari-proxy/v5
```

## Client library

`ari-proxy` uses semantic versioning and standard Go modules.  To use it in your
own Go package, simply reference the
`github.com/CyCoreSystems/ari-proxy/client/v5` package, and your dependency
management tool should be able to manage it.

### Usage

Connecting the client to NATS is simple:

```go
import (
   "github.com/CyCoreSystems/ari/v5"
   "github.com/CyCoreSystems/ari-proxy/v5/client"
)

func connect(ctx context.Context, appName string) (ari.Client,error) {
   c, err := client.New(ctx,
      client.WithApplication(appName),
      client.WithURI("nats://natshost:4222"),
   )
}
```

Connecting the client to RabbitMQ is like:

```go
import (
   "github.com/CyCoreSystems/ari/v5"
   "github.com/CyCoreSystems/ari-proxy/v5/client"
)

func connect(ctx context.Context, appName string) (ari.Client,error) {
   c, err := client.New(ctx,
      client.WithApplication(appName),
      client.WithURI("amqp://user:password@rabbitmqhost:5679/"),
   )
}
```

Configuration of the client can also be done with environment variables.
`ARI_APPLICATION` can be used to set the ARI application name, and `MESSAGEBUS_URL`
can be used to set the message bus URL.  Doing so allows you to get a client connection
simply with `client.New(ctx)`.

Once an `ari.Client` is obtained, the client functions exactly as the native
[ari](https://github.com/CyCoreSystems/ari) client.

More documentation:

  * [ARI library docs](https://godoc.org/github.com/CyCoreSystems/ari)

  * [ARI client examples](https://github.com/CyCoreSystems/ari/tree/master/_examples)


### Context

Note the use of the `context.Context` parameter.  This can be useful to properly
shutdown the client by a controlling context.  This shutdown will also close any
open subscriptions on the client.

Layers of clients can be used efficiently with different contexts using the
`New(context.Context)` function of each client instance.  Subtended clients will
be closed with their parents, use a common internal message bus connection, and can be
severally closed by their individual contexts.  This makes managing many active
channels easy and efficient.

### Lifecycle

There are two levels of client in use.  The first is a connection, which is a
long-lived network connection to the message bus.  In general, the end user
should not close this connection.  However, it is available, if necessary, as
`DefaultConn` and offers a `Close()` function for itself.

The second level is the ARI client.  Any number of ARI clients may make use of
the same underlying connection, but each client maintains its own separate bus
and subscription implementation.  Thus, when a user closes its client, the
connection is maintained, but all subscriptions are released.  Users should
always `Close()` their clients when done with them to avoid accumulating stale
subscriptions.

### Clustering

The ARI proxy works in a cluster setting by utilizing two coordinates:

 - The Asterisk ID
 - The ARI Application

Between the two of these, we can uniquely address each ARI resource, regardless
of where the client is located.  These pieces of information are handled
transparently and internally by the ARI proxy and the ARI proxy client to route
commands and events where they should be sent.

### Message bus protocol details

The protocol details described below are only necessary to know if you do not use the
provided client and/or server.  By using both components in this repository, the
protocol details below are transparently handled for you.

#### Subject structure

The message bus subject prefix defaults to `ari.`, and all messages used by this proxy
will be prefixed by that term.

Next is added one of four resource classifications:

 - `event` - Messages from Asterisk to clients
 - `get` - Read-only requests from clients to Asterisk
 - `command` - Non-creation operational requests from clients to Asterisk
 - `create` - Creation-related requests from clients to Asterisk

After the resource, the ARI application is appended. 

Finally, the Asterisk ID will be added to the end.  Thus, the subject for an event for the
ARI application "test" from the Asterisk box with ID "00:01:02:03:04:05" would
look like:

`ari.event.test.00:01:02:03:04:05`

For a channel creation command to the same app and node:

`ari.create.test.00:01:02:03:04:05`

The Asterisk ID component of the subject is optional for commands.  If a command
does not include an Asterisk ID, any ARI proxy running the provided ARI
application may respond to the request.  (Thus, implicitly, each ARI proxy
service listens to both its Asterisk ID-specific command subject and its generic
ARI application command subject.  In fact, each ARI proxy listens to each of the
three levels.  A request to `ari.command` will result in all ARI proxies
responding.)

This setup allows for a variable generalization in the listeners by using
message bus
wildcard subscriptions.  For instance, if you want to receive all events for the
"test" application regardless from which Asterisk machine they come, you would subscribe to:

`ari.event.test.>` //NATS
`ari.event.test.#` //RabbitMQ

#### Dialogs

Events may be further classified by the arbitrary "dialog" ID.  If any command
specifies a Dialog ID in its metadata, the ARI proxy will further classify
events related to that dialog.  Relationships are defined by the entity type on
which the Dialog-infused command operates.

Dialog-related events are published on their own message bus subject tree,
`dialogevent`.  Thus dialogs abstract ARI application and Asterisk ID.  An event
for dialog "testme123" would be published to:

`ari.dialogevent.testme123`

Keep in mind that regardless of dialog associations, all events are _also_
published to their appropriate canonical message bus subjects.  Dialogs are intended as
a mechanism to:

  - reduce client message traffic load
  - transcend ARI Applications and/or Asterisk nodes while maintaining logical
    separation of events

#### Message delivery

The means of a delivery for a generically-routed message depends on the type of
message it is.

  - Events are always delivered to all listeners.
  - Read-only commands are delivered to all listeners.
  - Non-creation operation commands are delivered to all listeners.
  - Creation-related commands are delivered to only one listener.

Thus, for efficiency, it is always recommended to use as precise a subject line
as possible.

#### Node discovery

Each ARI proxy sends out a periodic ping announcing itself in the cluster.
Clients may aggregate these pings to construct an expected map of machines in
the cluster.  Knowing this map allows the client to optimize its all-listener
commands by cancelling the wait period if it receives responses from all nodes
before the timeout has elapsed.

ARI proxies listen to `ari.ping` and send announcements on `ari.announce`.  The
structure of the announcement is thus:

```json
{
   "asterisk": "00:10:20:30:40:50",
   "application": "test"
}
```

#### Payload structure

For most requests, payloads exactly match their ARI library values.  However,
treatment of handlers is slightly different.

Instead of a handler, an `Entity` or array of `Entity`s is returned.  This
response type contains the Metadata for the entity (ARI application, Asterisk
ID, and optionally Dialog) as well as the unique ID of the entity.
