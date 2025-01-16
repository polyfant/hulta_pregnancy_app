package websocket

// Message represents a WebSocket message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
