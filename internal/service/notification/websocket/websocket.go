package websocket

import "context"

// WebSocketBroadcaster defines the interface for broadcasting messages via WebSocket
type WebSocketBroadcaster interface {
    Broadcast(ctx context.Context, message interface{}) error
}
