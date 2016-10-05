package session

// AppStart is the event sent on the start of an application and the creation of a server side dialog
type AppStart struct {
	ServerID    string `json:"server"`
	DialogID    string `json:"dialog"`
	Application string `json:"application"`
	ChannelID   string `json:"channel"` // The channel from the stasis start event
}
