package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketBroadcaster interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, userID string) error
	BroadcastToUser(ctx context.Context, userID string, notification *Notification) error
	Close()
}

type WebSocketBroadcasterImpl struct {
	clients  map[string]*websocket.Conn
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

func NewWebSocketBroadcaster() *WebSocketBroadcasterImpl {
	return &WebSocketBroadcasterImpl{
		clients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (wb *WebSocketBroadcasterImpl) HandleConnection(w http.ResponseWriter, r *http.Request, userID string) error {
	conn, err := wb.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("websocket upgrade error: %w", err)
	}

	wb.mu.Lock()
	wb.clients[userID] = conn
	wb.mu.Unlock()

	go wb.handleClientMessages(userID, conn)
	return nil
}

func (wb *WebSocketBroadcasterImpl) handleClientMessages(userID string, conn *websocket.Conn) {
	defer func() {
		wb.mu.Lock()
		delete(wb.clients, userID)
		conn.Close()
		wb.mu.Unlock()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error for user %s: %v", userID, err)
			break
		}
	}
}

func (wb *WebSocketBroadcasterImpl) BroadcastToUser(ctx context.Context, userID string, notification *Notification) error {
	wb.mu.RLock()
	defer wb.mu.RUnlock()

	client, exists := wb.clients[userID]
	if !exists {
		return fmt.Errorf("no websocket connection found for user %s", userID)
	}

	// Convert notification to JSON
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send notification via WebSocket
	if err := client.WriteMessage(websocket.TextMessage, notificationJSON); err != nil {
		log.Printf("Error broadcasting to user %s: %v", userID, err)
		return fmt.Errorf("failed to broadcast notification: %w", err)
	}

	return nil
}

func (wb *WebSocketBroadcasterImpl) Close() {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	for userID, conn := range wb.clients {
		conn.Close()
		delete(wb.clients, userID)
	}
}

// Ensure WebSocketBroadcaster implements the WebSocketBroadcaster interface
var _ WebSocketBroadcaster = (*WebSocketBroadcasterImpl)(nil)