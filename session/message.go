package session

// Message is the wrapper for a command sent over a dialog
type Message struct {
	Command string `json:"command"`
	Object  string `json:"object"`
	Payload []byte `json:"payload"`
}
