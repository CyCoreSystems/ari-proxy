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

	ChannelAnswer      *ChannelAnswer
	ChannelBusy        *ChannelBusy
	ChannelCongestion  *ChannelCongestion
	ChannelCreate      *ChannelCreate
	ChannelData        *ChannelData
	ChannelContinue    *ChannelContinue
	ChannelDial        *ChannelDial
	ChannelHangup      *ChannelHangup
	ChannelHold        *ChannelHold
	ChannelList        *ChannelList
	ChannelMOH         *ChannelMOH
	ChannelMute        *ChannelMute
	ChannelOriginate   *ChannelOriginate
	ChannelPlay        *ChannelPlay
	ChannelRecord      *ChannelRecord
	ChannelRing        *ChannelRing
	ChannelSendDTMF    *ChannelSendDTMF
	ChannelSilence     *ChannelSilence
	ChannelSnoop       *ChannelSnoop
	ChannelStopHold    *ChannelStopHold
	ChannelStopMOH     *ChannelStopMOH
	ChannelStopRing    *ChannelStopRing
	ChannelStopSilence *ChannelStopSilence
	ChannelSubscribe   *ChannelSubscribe
	ChannelUnmute      *ChannelUnmute
	ChannelVariables   *ChannelVariables

	DeviceStateData   *DeviceStateData
	DeviceStateDelete *DeviceStateDelete
	DeviceStateList   *DeviceStateList
	DeviceStateUpdate *DeviceStateUpdate

	/*
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
	Get *VariablesGet

	// Set is the Set variable request
	Set *VariablesSet
}

// VariablesGet is the request type for getting an asterisk variable
type VariablesGet struct {
	// VariablesGet is the signature field for this request
	VariablesGet struct{}
}

// VariablesSet is the request type for setting an asterisk variable
type VariablesSet struct {
	// VariablesSet is the signature field for this request
	VariablesSet struct{}

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

// BridgeSubscribe describes a request to subscribe a bridge
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

// ChannelCreate describes a request to create a new channel
type ChannelCreate struct {
	// ChannelCreate is the signature field for the request
	ChannelCreate struct{}

	// ChannelCreateRequest is the request for creating the channel
	ChannelCreateRequest ari.ChannelCreateRequest
}

// ChannelData describes a request to get the channel data
type ChannelData struct {
	// ChannelData is the signature field for the request
	ChannelData struct{}

	// ID is the channel ID
	ID string
}

// ChannelContinue describes a request to continue an ARI application
type ChannelContinue struct {
	// ChannelContinue is the signature field for the request
	ChannelContinue struct{}

	// ID is the channel ID
	ID string

	// ContinueRequest is the information for the continue request
	ContinueRequest client.ContinueRequest
}

// ChannelDial describes a request to dial
type ChannelDial struct {
	// ChannelDial is the signature field for the request
	ChannelDial struct{}

	// ID is the channel ID
	ID string

	// DialRequest is the data for the dial operation
	DialRequest client.DialRequest
}

// ChannelHangup is the request for hanging up a channel
type ChannelHangup struct {
	// ChannelHangup is the signature type for this request
	ChannelHangup struct{}

	// ID is the identifier for the channel
	ID string

	// Reason is the reason the channel is being hung up
	Reason string
}

// ChannelHold is the request for putting a channel on hold
type ChannelHold struct {

	// ChannelHold is the signature type for this request
	ChannelHold struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelList is the request for listing a channel
type ChannelList struct {
	// ChannelList is the signature type for this request
	ChannelList struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelMOH is the request playing hold on music on a channel
type ChannelMOH struct {
	// ChannelMOH is the signature type for this request
	ChannelMOH struct{}

	// ID is the identifier for the channel
	ID string

	// Music is the music to play
	Music string
}

// ChannelMute is the request for muting a channel
type ChannelMute struct {
	// ChannelMute is the signature type for this request
	ChannelMute struct{}

	// ID is the identifier for the channel
	ID string

	// Direction is the direction to mute
	Direction string
}

// ChannelOriginate is the request for creating a channel
type ChannelOriginate struct {
	// ChannelOriginate is the signature type for this request
	ChannelOriginate struct{}

	// OriginateRequest contains the information for originating a channel
	OriginateRequest ari.OriginateRequest
}

// ChannelPlay is the request for playing audio on a channel
type ChannelPlay struct {
	// ChannelPlay is the signature type for this request
	ChannelPlay struct{}

	// ID is the identifier for the channel
	ID string

	// PlayRequest is the request information for playing audio
	PlayRequest client.PlayRequest
}

// ChannelRecord is the request for recording a channel
type ChannelRecord struct {
	// ChannelRecord is the signature type for this request
	ChannelRecord struct{}

	// ID is the identifier for the channel
	ID string

	// RecordRequest is the recording request data
	RecordRequest client.RecordRequest
}

// ChannelRing is the request for playing a ringing noise on a channel
type ChannelRing struct {
	// ChannelRing is the signature type for this request
	ChannelRing struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelSendDTMF is the request for sending a DTMF event to a channel
type ChannelSendDTMF struct {
	// ChannelSendDTMF is the signature type for this request
	ChannelSendDTMF struct{}

	// ID is the identifier for the channel
	ID string

	// DTMF is the series of DTMF inputs to send
	DTMF string

	// Options are the DTMF options
	Options *ari.DTMFOptions
}

// ChannelSilence is the request for playing(?) silence on a channel
type ChannelSilence struct {
	// ChannelSilence is the signature type for this request
	ChannelSilence struct{}
	// ID is the identifier for the channel
	ID string
}

// ChannelSnoop is the request for snooping on a channel
type ChannelSnoop struct {
	// ChannelSnoop is the signature type for this request
	ChannelSnoop struct{}

	// ID is the identifier for the channel
	ID string

	// SnoopRequest is the request information for the snoop
	SnoopRequest client.SnoopRequest
}

// ChannelStopHold is the request for stopping the hold of a channel
type ChannelStopHold struct {
	// ChannelStopHold is the signature type for this request
	ChannelStopHold struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelStopMOH stops the music on old for a channel
type ChannelStopMOH struct {
	// ChannelStopMOH is the signature type for this request
	ChannelStopMOH struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelStopRing stops the ringing state for a channel
type ChannelStopRing struct {
	// ChannelStopRing is the signature type for this request
	ChannelStopRing struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelStopSilence stops the silence on the channel
type ChannelStopSilence struct {
	// ChannelStopSilence is the signature type for this request
	ChannelStopSilence struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelSubscribe describes the request for subscribing a channel to a dialog
type ChannelSubscribe struct {
	// ChannelSubscribe is the signature type for this request
	ChannelSubscribe struct{}

	// ID is the identifier for the channel
	ID string
}

// ChannelUnmute describes the request for unmuting the channel
type ChannelUnmute struct {
	// ChannelUnmute is the signature type for this request
	ChannelUnmute struct{}

	// ID is the identifier for the channel
	ID string

	// Direction is the direction of the unmute
	Direction string
}

// ChannelVariables is the request type for channel variable operations
type ChannelVariables struct {
	// ChannelVariables is the signature field for this request
	ChannelVariables struct{}

	// Name is the name of the channel
	ID string

	// Name is the name of the variable
	Name string

	// Get is the Get variable request
	Get *VariablesGet

	// Set is the Set variable request
	Set *VariablesSet
}

// DeviceStateData describes the request for getting the device state data
type DeviceStateData struct {
	// DeviceStateData is the signature type for this request
	DeviceStateData struct{}

	// ID is the identifier for the device
	ID string
}

// DeviceStateDelete describes the request for delete the device state
type DeviceStateDelete struct {
	// DeviceStateDelete is the signature type for this request
	DeviceStateDelete struct{}

	// ID is the identifier for the device
	ID string
}

// DeviceStateList describes the request for listing the devices and their states
type DeviceStateList struct {
	// DeviceStateList is the signature type for this request
	DeviceStateList struct{}

	// ID is the identifier for the device
	ID string
}

// DeviceStateUpdate describes the request for updating the device state
type DeviceStateUpdate struct {
	// DeviceStateUpdate is the signature type for this request
	DeviceStateUpdate struct{}

	// ID is the identifier for the device
	ID string

	// State is the new state of the device to set
	State string
}
