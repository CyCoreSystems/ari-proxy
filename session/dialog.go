package session

// A Dialog is a session between the ARI proxy client and the ARI proxy server
type Dialog struct {
	ID        string
	Transport Transport
	Objects   Objects
}

// NewDialog creates a new dialog with the given transport
func NewDialog(id string, transport Transport) *Dialog {
	return &Dialog{
		ID:        id,
		Transport: transport,
	}
}
