package proxy

import (
	"errors"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/client"
)

// AnnouncementInterval is the amount of time to wait between periodic service availability announcements
var AnnouncementInterval = time.Minute

// Announcement describes the structure of an ARI proxy's announcement of availability on the network.  These are sent periodically and upon request (by a Ping).
type Announcement struct {
	// Asterisk indicates the Asterisk ID to which the proxy is connected
	Asterisk string `json:"asterisk"`

	// Application indicates the ARI application as which the proxy is connected
	Application string `json:"application"`
}

// Metadata describes the metadata and associations of a message
type Metadata struct {
	// Application describes the ARI application
	Application string `json:"application,omitempty"`

	// Asterisk describes the ID of the associated Asterisk instance
	Asterisk string `json:"asterisk,omitempty"`

	// Dialog describes the dialog, if present
	Dialog string `json:"dialog,omitempty"`
}

// Entity is a response which returns a specific Entity, which is a stand-in for an entity Handler, containing the necessary descriptions to uniquely control the described entity.
type Entity struct {
	Metadata *Metadata

	// Type is the type of entity (application, asterisk, bridge, channel, deviceState, endpoint, mailbox, playback, liveRecording, storedRecording, sound)
	Type string `json:"type"`

	// ID is the unique identifier for the entity
	ID string `json:"name"`
}

// EntityList is a response which returns a list of Entities, as described above.
type EntityList struct {
	// List is the list of entities
	List []*Entity
}

// ErrNotFound indicates that the operation did not return a result
var ErrNotFound = errors.New("Not found")

// ErrorResponse is a response sent when a request could not be processed;
//
// NOTE: this is not always a problem, and sometimes it is expected (such as for
// a broadcast request for a particular channel, where only one proxy will have
// details for it).
type ErrorResponse struct {
	Error error
}

// NewErrorResponse wraps an error as an ErrorResponse
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Error: err}
}

// Event describes an ARI event sent from the ARI proxy to any subscribed clients
type Event struct {
	// Metadata is the metadata associated with the event
	Metadata *Metadata

	Event ari.Event
}

// Request describes a request which is sent from an ARI proxy Client to an ARI proxy Server
type Request struct {
	// Metadata is the metadata related to the request
	Metadata *Metadata

	ApplicationData        *ApplicationData
	ApplicationGet         *ApplicationGet
	ApplicationList        *ApplicationList
	ApplicationSubscribe   *ApplicationSubscribe
	ApplicationUnsubscribe *ApplicationUnsubscribe

	AsteriskInfo         *AsteriskInfo
	AsteriskReloadModule *AsteriskReloadModule
	AsteriskVariables    *AsteriskVariables
	//AsteriskConfig       *AsteriskConfig
	//AsteriskLogging      *AsteriskLogging
	//AsteriskModules      *AsteriskModules

	BridgeAddChannel    *BridgeAddChannel
	BridgeCreate        *BridgeCreate
	BridgeData          *BridgeData
	BridgeDelete        *BridgeDelete
	BridgeList          *BridgeList
	BridgePlay          *BridgePlay
	BridgeRecord        *BridgeRecord
	BridgeRemoveChannel *BridgeRemoveChannel
	BridgeSubscribe     *BridgeSubscribe

	ChannelAnswer     *ChannelAnswer
	ChannelBusy       *ChannelBusy
	ChannelCongestion *ChannelCongestion
	/*
		ChannelCreate       *ChannelCreate
		ChannelDataContinue *ChannelDataContinue
		ChannelDial         *ChannelDial
		ChannelHangup       *ChannelHangup
		ChannelHold         *ChannelHold
		ChannelList         *ChannelList
		ChannelMOH          *ChannelMOH
		ChannelMute         *ChannelMute
		ChannelOriginate    *ChannelOriginate
		ChannelPlay         *ChannelPlay
		ChannelRecord       *ChannelRecord
		ChannelRing         *ChannelRing
		ChannelSendDTMF     *ChannelSendDTMF
		ChannelSilence      *ChannelSilence
		ChannelSnoop        *ChannelSnoop
		ChannelStopHold     *ChannelStopHold
		ChannelStopMOH      *ChannelStopMOH
		ChannelStopRing     *ChannelStopRing
		ChannelStopSilence  *ChannelStopSilence
		ChannelSubscribe    *ChannelSubscribe
		ChannelUnmute       *ChannelUnmute
		ChannelVariables    *ChannelVariables

		DeviceStateData   *DeviceStateData
		DeviceStateDelete *DeviceStateDelete
		DeviceStateList   *DeviceStateList
		DeviceStateUpdate *DeviceStateUpdate

		EndpointData       *EndpointData
		EndpointGet        *EndpointGet
		EndpointList       *EndpointList
		EndpointListByTech *EndpointListByTech

		MailboxData   *MailboxData
		MailboxDelete *MailboxDelete
		MailboxList   *MailboxList
		MailboxUpdate *MailboxUpdate

		PlaybackControl   *PlaybackControl
		PlaybackData      *PlaybackData
		PlaybackStop      *PlaybackStop
		PlaybackSubscribe *PlaybackSubscribe

		RecordingStoredCopy   *RecordingStoredCopy
		RecordingStoredData   *RecordingStoredData
		RecordingStoredDelete *RecordingStoredDelete
		RecordingStoredList   *RecordingStoredList

		RecordingLiveData   *RecordingLiveData
		RecordingLiveDelete *RecordingLiveDelete
		RecordingLiveMute   *RecordingLiveMute
		RecordingLivePause  *RecordingLivePause
		RecordingLiveResume *RecordingLiveResume
		RecordingLiveScrap  *RecordingLiveScrap
		RecordingLiveStop   *RecordingLiveStop
		RecordingLiveUnmute *RecordingLiveUnmute

		SoundData *SoundData
		SoundList *SoundList
	*/
}

// ApplicationData describes a request to get the data for a particular ARI application
type ApplicationData struct {
	// ApplicationData is the signature field for this request
	ApplicationData struct{}

	// Name is the name of the ARI application to be retrieved
	Name string
}

// ApplicationGet describes a request for a particular ARI application
type ApplicationGet struct {
	// ApplicationGet is the signature field for this request
	ApplicationGet struct{}

	// Name is the name of the ARI application to be retrieved
	Name string
}

// ApplicationList describes a request for the list of ARI applications
type ApplicationList struct {
	// ApplicationList is the signature field for this request
	ApplicationList struct{}
}

// ApplicationSubscribe describes a request to subscribe a particular ARI application to an EventSource
type ApplicationSubscribe struct {
	// ApplicationSubscribe is the signature field for this request
	ApplicationSubscribe struct{}

	// Name is the name of the ARI application to be retrieved
	Name string

	// EventSource is the ARI event source to which the subscription is requested.  This should be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	EventSource string
}

// ApplicationUnsubscribe describes a request to unsubscribe a particular ARI application from an EventSource
type ApplicationUnsubscribe struct {
	// ApplicationUnsubscribe is the signature field for this request
	ApplicationUnsubscribe struct{}

	// Name is the name of the ARI application to be retrieved
	Name string

	// EventSource is the ARI event source of which the unsubscription is requested.  This should be one of:
	//  - channel:<channelId>
	//  - bridge:<bridgeId>
	//  - endpoint:<tech>/<resource> (e.g. SIP/102)
	//  - deviceState:<deviceName>
	EventSource string
}

// AsteriskInfo describes a request to get the asterisk information
type AsteriskInfo struct {
	// AsteriskInfo is the signature field for this request
	AsteriskInfo struct{}
}

// AsteriskReloadModule descibres a request to reload an asterisk module
type AsteriskReloadModule struct {
	// AsteriskReloadModule is the signature field for this request
	AsteriskReloadModule struct{}

	// Name is the name of the asterisk module to reload
	Name string
}

// AsteriskVariables is the request type for asterisk variable operations
type AsteriskVariables struct {
	// AsteriskVariables is the signature field for this request
	AsteriskVariables struct{}

	// Name is the name of the asterisk variable
	Name string

	// Get is the Get variable request
	Get *AsteriskVariablesGet

	// Set is the Set variable request
	Set *AsteriskVariablesSet
}

// AsteriskVariablesGet is the request type for getting an asterisk variable
type AsteriskVariablesGet struct {
	// AsteriskVariablesGet is the signature field for this request
	AsteriskVariablesGet struct{}
}

// AsteriskVariablesSet is the request type for setting an asterisk variable
type AsteriskVariablesSet struct {
	// AsteriskVariablesSet is the signature field for this request
	AsteriskVariablesGet struct{}

	// Value is the value to set
	Value string
}

// BridgeAddChannel is the request type for adding a channel to a bridge
type BridgeAddChannel struct {
	// BridgeAddChannel is the signature field for this request
	BridgeAddChannel struct{}

	// Name is the name of the bridge
	ID string

	// Channel is the channel to add to the bridge
	Channel string
}

// BridgeCreate is the request type for creating a bridge
type BridgeCreate struct {
	// BridgeCreate is the signature field for the request
	BridgeCreate struct{}

	// ID is the id of the bridge
	CreateBridgeRequest client.CreateBridgeRequest
}

// BridgeData is the request type for getting the bridge data
type BridgeData struct {
	// BridgeData is the signature field for the request
	BridgeData struct{}

	// ID is the identifier of the bridge to get
	ID string
}

// BridgeDelete is the request type for deleting a bridge
type BridgeDelete struct {
	// BridgeDelete is the signature field for the request
	BridgeDelete struct{}

	// ID is the identifier of the bridge
	ID string
}

// BridgeList is the request type for listing the bridges
type BridgeList struct {
	// BridgeList is the signature field for the request
	BridgeList struct{}
}

// BridgePlay is the request type for playing audio on the bridge
type BridgePlay struct {
	// BridgePlay is the signature field for the request
	BridgePlay struct{}

	// ID is the identifier of the bridge
	ID string

	// PlayRequest is the request for the playing of audio
	PlayRequest client.PlayRequest
}

// BridgeRecord is the request for recording a bridge
type BridgeRecord struct {
	// BridgeRecord is the signature field for this request
	BridgeRecord struct{}

	// ID is the identifier of the bridge
	ID string

	// RecordRequest is the request for recording audio
	RecordRequest client.RecordRequest
}

// BridgeRemoveChannel is the request for removing a channel on the bridge
type BridgeRemoveChannel struct {
	// BridgeRemoveChannel is the signature field for this request
	BridgeRemoveChannel struct{}

	// ID is the identifier of the bridge
	ID string

	// Channel is the name of the channel to remove
	Channel string
}

// BridgeSubscribe describes a request to subscribe a particular ARI application to an EventSource
type BridgeSubscribe struct {
	// ApplicationSubscribe is the signature field for this request
	BridgeSubscribe struct{}

	// ID is the identifier of the bridge
	ID string
}

// ChannelAnswer describes a request to answer a channel
type ChannelAnswer struct {
	// ChannelAnswer is the signature field for the request
	ChannelAnswer struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelBusy describes a request to send a busy signal to a channel
type ChannelBusy struct {
	// ChannelBusy is the signature field for the request
	ChannelBusy struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelCongestion describes a request to send a congestion signal to a channel
type ChannelCongestion struct {
	// ChannelCongestion is the signature field for the request
	ChannelCongestion struct{}

	// ID is the identifier for the channel
	ID string
}
