# ari-proxy

Proxy for the Asterisk REST interface.

Version 1 of this system sought to maintain compatibility with the other ARI
proxies.  Version 2 is written to address many shortcomings of the Version 1
system and fundamentally alters the wire protocol.  End-user API, however,
should be identical.

## Installation

`master` should be the latest stable, so a simple `go get` is required:

	go get github.com/CyCoreSystems/ari-proxy/cmd/ari-proxy

## Development

New development occurs on `develop` branch.

## Client lifecycle

There are two levels of client in use.  The first is a connection, which is a
long-lived network connection to the NATS cluster.  In general, the end user
should not close this connection.  However, it is available, if necessary, as
`DefaultConn` and offers a `Close()` function for itself.

The second level is the ARI client.  Any number of ARI clients may make use of
the same underlying connection, but each client maintains its own separate bus
and subscription implementation.  Thus, when a user closes its client, the
connection is maintained, but all subscriptions are released.  Users should
always `Close()` their clients when done with them to avoid accumulating stale
subscriptions.

## Clustering

The ARI proxy works in a cluster setting by utilizing two coordinates:

 - The Asterisk ID
 - The ARI Application

Between the two of these, we can uniquely address each ARI resource, regardless
of where the client is located.  These pieces of information are handled
transparently and internally by the ARI proxy and the ARI proxy client to route
commands and events where they should be sent.

## NATS protocol details

The protocol details described below are only necessary to know if you do not use the
provided client and/or server.  By using both components in this repository, the
protocol details below are transparently handled for you.

### Subject structure

The NATS subject prefix defaults to `ari.`, and all messages used by this proxy
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

This setup allows for a variable generalization in the listeners by using NATS
wildcard subscriptions.  For instance, if you want to receive all events for the
"test" application regardless from which Asterisk machine they come, you would
subscribe to:

`ari.event.test.>`

### Dialogs

Events may be further classified by the arbitrary "dialog" ID.  If any command
specifies a Dialog ID in its metadata, the ARI proxy will further classify
events related to that dialog.  Relationships are defined by the entity type on
which the Dialog-infused command operates.

Dialog-related events are published on their own NATS subject tree,
`dialogevent`.  Thus dialogs abstract ARI application and Asterisk ID.  An event
for dialog "testme123" would be published to:

`ari.dialogevent.testme123`

Keep in mind that regardless of dialog associations, all events are _also_
published to their appropriate canonical NATS subjects.  Dialogs are intended as
a mechanism to:

  - reduce client message traffic load
  - transcend ARI Applications and/or Asterisk nodes while maintaining logical
    separation of events

### Message delivery

The means of a delivery for a generically-routed message depends on the type of
message it is.

  - Events are always delivered to all listeners.
  - Read-only commands are delivered to all listeners.
  - Non-creation operation commands are delivered to all listeners.
  - Creation-related commands are delivered to only one listener.

Thus, for efficiency, it is always recommended to use as precise a subject line
as possible.

### Node discovery

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


