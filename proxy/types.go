package proxy

import (
	"errors"
	"fmt"
	"time"

	"github.com/CyCoreSystems/ari"
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

// AnnouncementSubject returns the NATS subject
func AnnouncementSubject(prefix string) string {
	return fmt.Sprintf("%sannounce", prefix)
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
	Metadata *Metadata `json:"metadata"`

	// Kind is the type of entity (application, asterisk, bridge, channel, deviceState, endpoint, logging, mailbox, playback, liveRecording, storedRecording, sound)
	Kind string `json:"kind"`

	// ID is the unique identifier for the entity
	ID string `json:"name"`
}

// Entity types
const (
	// EntityKindApplication indicates that the entity is an ARI Application
	EntityKindApplication = "application"

	// EntityKindAsterisk indicates that the entity is an Asterisk
	EntityKindAsterisk = "asterisk"

	// EntityKindBridge indicates that the entity is a Bridge
	EntityKindBridge = "bridge"

	// EntityKindChannel indicates that the entity is a Channel
	EntityKindChannel = "channel"

	// EntityKindDeviceState indicates that the entity is a DeviceState
	EntityKindDeviceState = "device_state"

	// EntityKindEndpoint indicates that the entity is an Endpoint
	EntityKindEndpoint = "endpoint"

	// EntityKindLogging indicates that the entity is a Logging channel
	EntityKindLogging = "logging"

	// EntityKindMailbox indicates that the entity is a Mailbox
	EntityKindMailbox = "mailbox"

	// EntityKindPlayback indicates that the entity is a Playback
	EntityKindPlayback = "playback"

	// EntityKindLiveRecording indicates that the entity is a Live Recording
	EntityKindLiveRecording = "live_recording"

	// EntityKindStoredRecording indicates that the entity is a Stored Recording
	EntityKindStoredRecording = "stored_recording"

	// EntityKindSound indicates that the entity is a Sound
	EntityKindSound = "sound"
)

// EntityData is a response which returns the data for a specific entity.
type EntityData struct {
	Metadata *Metadata `json:"metadata"`

	Application     *ari.ApplicationData     `json:"application,omitempty"`
	Asterisk        *ari.AsteriskInfo        `json:"asterisk,omitempty"`
	Bridge          *ari.BridgeData          `json:"bridge,omitempty"`
	Channel         *ari.ChannelData         `json:"channel,omitempty"`
	Config          *ari.ConfigData          `json:"config,omitempty"`
	DeviceState     *ari.DeviceStateData     `json:"device_state,omitempty"`
	Endpoint        *ari.EndpointData        `json:"endpoint,omitempty"`
	LiveRecording   *ari.LiveRecordingData   `json:"live_recording,omitempty"`
	Log             *ari.LogData             `json:"log,omitempty"`
	Mailbox         *ari.MailboxData         `json:"mailbox,omitempty"`
	Module          *ari.ModuleData          `json:"module,omitempty"`
	Playback        *ari.PlaybackData        `json:"playback,omitempty"`
	Sound           *ari.SoundData           `json:"sound,omitempty"`
	StoredRecording *ari.StoredRecordingData `json:"stored_recording,omitempty"`
	TextMessage     *ari.TextMessageData     `json:"text_message,omitempty"`

	Variable string `json:"variable,omitempty"`
}

// EntityList is a response which returns a list of Entities, as described above.
type EntityList struct {
	// List is the list of entities
	List []*Entity `json:"list,omitempty"`
}

// ErrNotFound indicates that the operation did not return a result
var ErrNotFound = errors.New("Not found")

// Response is a response to a request.  This acts as a base type for more complicated responses, as well.
type Response struct {
	Error string `json:"error,omitempty"`

	// Data is the returned entity data, if applicable
	Data *EntityData `json:"data,omitempty"`

	// Entity is the returned entity, if applicable
	Entity *Entity `json:"entity,omitempty"`

	// EntityList is the returned list of entities, if applicable
	EntityList *EntityList `json:"list,omitempty"`
}

// Err returns an error from the Response.  If the response's Error is empty, a nil error is returned.  Otherwise, the error will be filled with the value of response.Error.
func (e *Response) Err() error {
	if e.Error != "" {
		return errors.New(e.Error)
	}
	return nil
}

// IsNotFound indicates that the retuned error response was a Not Found error response
func (e *Response) IsNotFound() bool {
	return e.Error == "Not found"
}

// NewErrorResponse wraps an error as an ErrorResponse
func NewErrorResponse(err error) *Response {
	if err == nil {
		return &Response{}
	}
	return &Response{Error: err.Error()}
}

// Request describes a request which is sent from an ARI proxy Client to an ARI proxy Server
type Request struct {
	// Metadata is the metadata related to the request
	Metadata *Metadata `json:"metadata"`

	ApplicationData        *ApplicationData        `json:"application_data,omitempty"`
	ApplicationGet         *ApplicationGet         `json:"application_get,omitempty"`
	ApplicationList        *ApplicationList        `json:"application_list,omitempty"`
	ApplicationSubscribe   *ApplicationSubscribe   `json:"application_subscribe,omitempty"`
	ApplicationUnsubscribe *ApplicationUnsubscribe `json:"application_unsubscribe,omitempty"`

	AsteriskInfo         *AsteriskInfo         `json:"asterisk_info,omitempty"`
	AsteriskReloadModule *AsteriskReloadModule `json:"asterisk_reload_module,omitempty"`
	AsteriskVariables    *AsteriskVariables    `json:"asterisk_variables,omitempty"`
	AsteriskConfig       *AsteriskConfig       `json:"asterisk_config,omitempty"`
	AsteriskLogging      *AsteriskLogging      `json:"asterisk_logging,omitempty"`
	AsteriskModules      *AsteriskModules      `json:"asterisk_modules,omitempty"`

	BridgeAddChannel    *BridgeAddChannel    `json:"bridge_add_channel,omitempty"`
	BridgeCreate        *BridgeCreate        `json:"bridge_create,omitempty"`
	BridgeData          *BridgeData          `json:"bridge_data,omitempty"`
	BridgeDelete        *BridgeDelete        `json:"bridge_delete,omitempty"`
	BridgeList          *BridgeList          `json:"bridge_list,omitempty"`
	BridgePlay          *BridgePlay          `json:"bridge_play,omitempty"`
	BridgeRecord        *BridgeRecord        `json:"bridge_record,omitempty"`
	BridgeRemoveChannel *BridgeRemoveChannel `json:"bridge_remove_channel,omitempty"`
	BridgeSubscribe     *BridgeSubscribe     `json:"bridge_subscribe,omitempty"`

	ChannelAnswer      *ChannelAnswer      `json:"channel_answer,omitempty"`
	ChannelBusy        *ChannelBusy        `json:"channel_busy,omitempty"`
	ChannelCongestion  *ChannelCongestion  `json:"channel_congestion,omitempty"`
	ChannelCreate      *ChannelCreate      `json:"channel_create,omitempty"`
	ChannelData        *ChannelData        `json:"channel_data,omitempty"`
	ChannelContinue    *ChannelContinue    `json:"channel_continue,omitempty"`
	ChannelDial        *ChannelDial        `json:"channel_dial,omitempty"`
	ChannelHangup      *ChannelHangup      `json:"channel_hangup,omitempty"`
	ChannelHold        *ChannelHold        `json:"channel_hold,omitempty"`
	ChannelList        *ChannelList        `json:"channel_list,omitempty"`
	ChannelMOH         *ChannelMOH         `json:"channel_moh,omitempty"`
	ChannelMute        *ChannelMute        `json:"channel_mute,omitempty"`
	ChannelOriginate   *ChannelOriginate   `json:"channel_originate,omitempty"`
	ChannelPlay        *ChannelPlay        `json:"channel_play,omitempty"`
	ChannelRecord      *ChannelRecord      `json:"channel_record,omitempty"`
	ChannelRing        *ChannelRing        `json:"channel_ring,omitempty"`
	ChannelSendDTMF    *ChannelSendDTMF    `json:"channel_send_dtmf,omitempty"`
	ChannelSilence     *ChannelSilence     `json:"channel_silence,omitempty"`
	ChannelSnoop       *ChannelSnoop       `json:"channel_snoop,omitempty"`
	ChannelStopHold    *ChannelStopHold    `json:"channel_stop_hold,omitempty"`
	ChannelStopMOH     *ChannelStopMOH     `json:"channel_stop_moh,omitempty"`
	ChannelStopRing    *ChannelStopRing    `json:"channel_stop_ring,omitempty"`
	ChannelStopSilence *ChannelStopSilence `json:"channel_stop_silence,omitempty"`
	ChannelSubscribe   *ChannelSubscribe   `json:"channel_subscribe,omitempty"`
	ChannelUnmute      *ChannelUnmute      `json:"channel_unmute,omitempty"`
	ChannelVariables   *ChannelVariables   `json:"channel_variables,omitempty"`

	DeviceStateData   *DeviceStateData   `json:"device_state_data,omitempty"`
	DeviceStateDelete *DeviceStateDelete `json:"device_state_delete,omitempty"`
	DeviceStateList   *DeviceStateList   `json:"device_state_list,omitempty"`
	DeviceStateUpdate *DeviceStateUpdate `json:"device_state_update,omitempty"`

	EndpointData       *EndpointData       `json:"endpoint_data,omitempty"`
	EndpointList       *EndpointList       `json:"endpoint_list,omitempty"`
	EndpointListByTech *EndpointListByTech `json:"endpoint_list_by_tech,omitempty"`

	MailboxData   *MailboxData   `json:"mailbox_data,omitempty"`
	MailboxDelete *MailboxDelete `json:"mailbox_delete,omitempty"`
	MailboxList   *MailboxList   `json:"mailbox_list,omitempty"`
	MailboxUpdate *MailboxUpdate `json:"mailbox_update,omitempty"`

	PlaybackControl   *PlaybackControl   `json:"playback_control,omitempty"`
	PlaybackData      *PlaybackData      `json:"playback_data,omitempty"`
	PlaybackStop      *PlaybackStop      `json:"playback_stop,omitempty"`
	PlaybackSubscribe *PlaybackSubscribe `json:"playback_subscribe,omitempty"`

	RecordingStoredCopy   *RecordingStoredCopy   `json:"recording_stored_copy,omitempty"`
	RecordingStoredData   *RecordingStoredData   `json:"recording_stored_data,omitempty"`
	RecordingStoredDelete *RecordingStoredDelete `json:"recording_stored_delete,omitempty"`
	RecordingStoredList   *RecordingStoredList   `json:"recording_stored_list,omitempty"`

	RecordingLiveData   *RecordingLiveData   `json:"recording_live_data,omitempty"`
	RecordingLiveDelete *RecordingLiveDelete `json:"recording_live_delete,omitempty"`
	RecordingLiveMute   *RecordingLiveMute   `json:"recording_live_mute,omitempty"`
	RecordingLivePause  *RecordingLivePause  `json:"recording_live_pause,omitempty"`
	RecordingLiveResume *RecordingLiveResume `json:"recording_live_resume,omitempty"`
	RecordingLiveScrap  *RecordingLiveScrap  `json:"recording_live_scrap,omitempty"`
	RecordingLiveStop   *RecordingLiveStop   `json:"recording_live_stop,omitempty"`
	RecordingLiveUnmute *RecordingLiveUnmute `json:"recording_live_unmute,omitempty"`

	SoundData *SoundData `json:"sound_data,omitempty"`
	SoundList *SoundList `json:"sound_list,omitempty"`
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
	ID string

	// Type is the comma-separated list of bridge type attributes (mixing, holding, dtmf_events, proxy_media)
	Type string

	// Name is the name to assign to the bridge (optional)
	Name string
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

	// PlaybackID is the unique identifier for this playback
	PlaybackID string

	// MediaURI is the URI from which to obtain the playback media
	MediaURI string
}

// BridgeRecord is the request for recording a bridge
type BridgeRecord struct {
	// BridgeRecord is the signature field for this request
	BridgeRecord struct{}

	// ID is the identifier of the bridge
	ID string

	// Name is the name for the recording
	Name string

	// Options is the list of recording Options
	Options *ari.RecordingOptions
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

	// Context is the context into which the channel should be continued
	Context string

	// Extension is the extension into which the channel should be continued
	Extension string

	// Priority is the priority at which the channel should be continued
	Priority int
}

// ChannelDial describes a request to dial
type ChannelDial struct {
	// ChannelDial is the signature field for the request
	ChannelDial struct{}

	// ID is the channel ID
	ID string

	// Caller is the channel ID of the "caller" channel; if specified, the media parameters of the dialing channel will be matched to the "caller" channel.
	Caller string

	// Timeout is the maximum time which should be allowed for the dial to complete
	Timeout time.Duration
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
	Direction ari.Direction
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

	// PlaybackID is the unique identifier for this playback
	PlaybackID string

	// MediaURI is the URI from which to obtain the playback media
	MediaURI string
}

// ChannelRecord is the request for recording a channel
type ChannelRecord struct {
	// ChannelRecord is the signature type for this request
	ChannelRecord struct{}

	// ID is the identifier for the channel
	ID string

	// Name is the name for the recording
	Name string

	// Options is the list of recording Options
	Options *ari.RecordingOptions
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

	// SnoopID is the ID to use for the snoop channel which will be created.
	SnoopID string

	// Options describe the parameters for the snoop session
	Options *ari.SnoopOptions
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
	Direction ari.Direction
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

// EndpointData describes the request for getting endpoint data
type EndpointData struct {
	// EndpointData is the signature type for this request
	EndpointData struct{}

	// Tech is the technology for the endpoint
	Tech string

	// Resource is the resource for the endpoint
	Resource string
}

// EndpointList describes the request for the listing endpoints
type EndpointList struct {
	// EndpointList is the signature type for this request
	EndpointList struct{}
}

// EndpointListByTech describes the request for listing endpoints by technology
type EndpointListByTech struct {
	// EndpointListByTech is the signature type for this request
	EndpointListByTech struct{}

	// Tech is the technology for the endpoint
	Tech string
}

// MailboxData describes the request for getting the mailbox data
type MailboxData struct {
	// MailboxData is the signature type for this request
	MailboxData struct{}

	// Name is the name of the mailbox
	Name string
}

// MailboxDelete describes the request for deleting a mailbox
type MailboxDelete struct {
	// MailboxDelete is the signature type for this request
	MailboxDelete struct{}

	// Name is the name of the mailbox
	Name string
}

// MailboxList describes the request for listing mailboxes
type MailboxList struct {
	// MailboxList is the signature type for this request
	MailboxList struct{}
}

// MailboxUpdate describes the request for updating a mailbox
type MailboxUpdate struct {
	// MailboxUpdate is the signature type for this request
	MailboxUpdate struct{}

	// Name is the name of the mailbox
	Name string

	// New is the number of New (unread) messages in the mailbox
	New int

	// Old is the number of Old (read) messages in the mailbox
	Old int
}

// PlaybackControl describes the request for performing a playback command
type PlaybackControl struct {
	// PlaybackControl is the signature type for this request
	PlaybackControl struct{}

	// ID is the playback identifier
	ID string

	// Command is the playback control command to run
	Command string
}

// PlaybackData describes the request for getting the playback data
type PlaybackData struct {
	// PlaybackData is the signature type for this request
	PlaybackData struct{}

	// ID is the playback identifier
	ID string
}

// PlaybackStop describes the request for stopping a playback
type PlaybackStop struct {
	// PlaybackStop is the signature type for this request
	PlaybackStop struct{}

	// ID is the playback identifier
	ID string
}

// PlaybackSubscribe describes the request for binding a playback object to the dialog
type PlaybackSubscribe struct {
	// PlaybackSubscribe is the signature type for this request
	PlaybackSubscribe struct{}

	// ID is the playback identifier
	ID string
}

// RecordingStoredCopy describes the request for copying a stored recording
type RecordingStoredCopy struct {
	// RecordingStoredCopy is the signature type for this request
	RecordingStoredCopy struct{}

	// ID is the stored recording identifier
	ID string

	// Destination is the destination location to copy to
	Destination string
}

// RecordingStoredData describes the request for getting the stored recording data
type RecordingStoredData struct {
	// RecordingStoredData is the signature type for this request
	RecordingStoredData struct{}

	// ID is the stored recording identifier
	ID string
}

// RecordingStoredDelete describes the request for deleting the stored recording
type RecordingStoredDelete struct {
	// RecordingStoredDelete is the signature type for this request
	RecordingStoredDelete struct{}

	// ID is the stored recording identifier
	ID string
}

// RecordingStoredList describes the request for listing the stored recordings
type RecordingStoredList struct {
	// RecordingStoredList is the signature type for this request
	RecordingStoredList struct{}

	// ID is the stored recording identifier
	ID string
}

// RecordingLiveData decribes the request for getting the recording data
type RecordingLiveData struct {
	// RecordingLiveData is the signature type for this request
	RecordingLiveData struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveDelete describes the request for deleting the live recording
type RecordingLiveDelete struct {
	// RecordingLiveDelete is the signature type for this request
	RecordingLiveDelete struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveMute describes the request for muting a live recording
type RecordingLiveMute struct {
	// RecordingLiveMute is the signature type for this request
	RecordingLiveMute struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLivePause describes the request for pausing a live recording
type RecordingLivePause struct {
	// RecordingLivePause is the signature type for this request
	RecordingLivePause struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveResume describes the request for resuming a live recording
type RecordingLiveResume struct {
	// RecordingLiveResume is the signature type for this request
	RecordingLiveResume struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveScrap describes the request for scrapping a live recording
type RecordingLiveScrap struct {
	// RecordingLiveScrap is the signature type for this request
	RecordingLiveScrap struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveStop describes the request for stopping a live recording
type RecordingLiveStop struct {
	// RecordingLiveStop is the signature type for this request
	RecordingLiveStop struct{}

	// ID is the live recording identifier
	ID string
}

// RecordingLiveUnmute describes the request for unmuting a live recording
type RecordingLiveUnmute struct {
	// RecordingLiveUnmute is the signature type for this request
	RecordingLiveUnmute struct{}

	// ID is the live recording identifier
	ID string
}

// SoundData describes the request for getting the sound data
type SoundData struct {
	// SoundData is the signature type for this request
	SoundData struct{}

	// Name is the name of the sound
	Name string
}

// SoundList describes the request for listing the sounds
type SoundList struct {
	// SoundList is the signature type for this request
	SoundList struct{}

	// Filters are the filters to apply when listing the sounds
	Filters map[string]string
}

// AsteriskConfig describes the requests for asterisk configuration
type AsteriskConfig struct {
	// AsteriskConfig is the signature type for this request
	AsteriskConfig struct{}

	// ConfigClass is the class of the configuration
	ConfigClass string

	// ObjectType is the type of the configuration object
	ObjectType string

	// ID is the configuration identifier
	ID string

	// Data is the asterisk config get data request
	Data *AsteriskConfigData

	// Delete is the asterisk delete config request
	Delete *AsteriskConfigDelete

	// Update is the asterisk update config request
	Update *AsteriskConfigUpdate
}

// AsteriskConfigData describes the request for getting asterisk configuration data
type AsteriskConfigData struct {
	// AsteriskConfigData is the signature type for this request
	AsteriskConfigData struct{}
}

// AsteriskConfigUpdate describes the request for updating asterisk configuration data
type AsteriskConfigUpdate struct {
	// AsteriskConfigUpdate is the signature type for this request
	AsteriskConfigUpdate struct{}

	// Tuples is the list of configuration tuples to update
	Tuples []ari.ConfigTuple
}

// AsteriskConfigDelete describes the request for deleting asterisk configuration data
type AsteriskConfigDelete struct {
	// AsteriskConfigDelete is the signature type for this request
	AsteriskConfigDelete struct{}
}

// AsteriskLogging describes the group of requests for asterisk logging operations
type AsteriskLogging struct {
	// AsteriskLogging is the signature type for this request
	AsteriskLogging struct{}

	// Create is the asterisk logging create request
	Create *AsteriskLoggingCreate

	// Data is a logging channel data request
	Data *AsteriskLoggingData

	// Delete is the asterisk logging delete request
	Delete *AsteriskLoggingDelete

	// List is the asterisk logging list request
	List *AsteriskLoggingList

	// Rotate is the aterisk logging rotate request
	Rotate *AsteriskLoggingRotate
}

// AsteriskLoggingList describes the asterisk logging list request
type AsteriskLoggingList struct {
	// AsteriskLoggingList is the signature type for this request
	AsteriskLoggingList struct{}
}

// AsteriskLoggingCreate describes the asterisk logging create request
type AsteriskLoggingCreate struct {
	// AsteriskLoggingCreate is the signature type for this request
	AsteriskLoggingCreate struct{}

	// ID is the identifier for this object
	ID string

	// Config is the config details for the logging object
	Config string
}

// AsteriskLoggingData describes the asterisk logging data request
type AsteriskLoggingData struct {
	// AsteriskLoggingData is the signature type for this request
	AsteriskLoggingData struct{}

	// ID is the identifier for this object
	ID string
}

// AsteriskLoggingDelete describes the asterisk logging delete request
type AsteriskLoggingDelete struct {
	// AsteriskLoggingDelete is the signature type for this request
	AsteriskLoggingDelete struct{}

	// ID is the identifier for this object
	ID string
}

// AsteriskLoggingRotate describes the asterisk logging rotate request
type AsteriskLoggingRotate struct {
	// AsteriskLoggingRotate is the signature type for this request
	AsteriskLoggingRotate struct{}

	// ID is the identifier for this object
	ID string
}

// AsteriskModules describes the group of operations on asterisk modules
type AsteriskModules struct {
	// AsteriskModules is the signature type for this request
	AsteriskModules struct{}

	// List is the asterisk modules list operation
	List *AsteriskModulesList

	// Data is the asterisk modules get data operation
	Data *AsteriskModulesData

	// Load is the asterisk modules load operation
	Load *AsteriskModulesLoad

	// Unload is the asterisk modules unload operation
	Unload *AsteriskModulesUnload

	// Reload is the asterisk modules unload operation
	Reload *AsteriskModulesReload
}

// AsteriskModulesList describes the asterisk modules list request
type AsteriskModulesList struct {
	// AsteriskModulesList is the signature type for this request
	AsteriskModulesList struct{}
}

// AsteriskModulesData describes the asterisk get data request
type AsteriskModulesData struct {
	// AsteriskModulesData is the signature type for this request
	AsteriskModulesData struct{}

	// Name is the name of the asterisk module
	Name string
}

// AsteriskModulesLoad describes the asterisk load module request
type AsteriskModulesLoad struct {
	// AsteriskModulesLoad is the signature type for this request
	AsteriskModulesLoad struct{}

	// Name is the name of the asterisk module
	Name string
}

// AsteriskModulesUnload describes the asterisk unload module request
type AsteriskModulesUnload struct {
	// AsteriskModulesUnload is the signature type for this request
	AsteriskModulesUnload struct{}

	// Name is the name of the asterisk module
	Name string
}

// AsteriskModulesReload describes the asterisk reload module request
type AsteriskModulesReload struct {
	// AsteriskModulesReload is the signature type for this request
	AsteriskModulesReload struct{}

	// Name is the name of the asterisk module
	Name string
}
