package service

import (
	"context"
	"net/http"

	"github.com/polyfant/hulta_pregnancy_app/internal/service/notification"
)

// WebSocketBroadcaster defines the interface for broadcasting notifications over WebSocket
type WebSocketBroadcaster interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, userID string) error
	BroadcastToUser(ctx context.Context, userID string, notification *notification.Notification) error
	Close()
}
