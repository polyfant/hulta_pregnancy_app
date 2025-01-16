package websocket

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocketHub(t *testing.T) {
	t.Run("Client Registration", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create test client
		client := &Client{
			hub:     hub,
			send:    make(chan []byte, 256),
			horseID: 1,
		}

		// Register client
		hub.register <- client

		// Wait for registration
		time.Sleep(100 * time.Millisecond)

		// Verify registration
		hub.mu.RLock()
		clients := hub.clients[client.horseID]
		hub.mu.RUnlock()

		assert.NotNil(t, clients)
		assert.True(t, clients[client])
	})

	t.Run("Client Unregistration", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create and register client
		client := &Client{
			hub:     hub,
			send:    make(chan []byte, 256),
			horseID: 1,
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)

		// Unregister client
		hub.unregister <- client
		time.Sleep(100 * time.Millisecond)

		// Verify unregistration
		hub.mu.RLock()
		clients := hub.clients[client.horseID]
		hub.mu.RUnlock()

		assert.Empty(t, clients)
	})

	t.Run("Message Broadcasting", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create test server
		server := testutil.CreateTestWebSocketServer(hub)
		defer server.Close()

		// Create test clients
		conn1, err := testutil.CreateTestWebSocketClient(server, "/ws/1")
		require.NoError(t, err)
		defer conn1.Close()

		conn2, err := testutil.CreateTestWebSocketClient(server, "/ws/1")
		require.NoError(t, err)
		defer conn2.Close()

		// Register clients
		client1 := &Client{
			hub:     hub,
			conn:    conn1,
			send:    make(chan []byte, 256),
			horseID: 1,
		}
		client2 := &Client{
			hub:     hub,
			conn:    conn2,
			send:    make(chan []byte, 256),
			horseID: 1,
		}

		hub.register <- client1
		hub.register <- client2
		time.Sleep(100 * time.Millisecond)

		// Broadcast test message
		testMessage := map[string]interface{}{
			"type": "TEST",
			"data": "Hello, World!",
		}
		hub.BroadcastToHorse(1, testMessage)

		// Verify both clients receive the message
		assert.NoError(t, testutil.AssertWebSocketMessage(t, conn1, time.Second, testMessage))
		assert.NoError(t, testutil.AssertWebSocketMessage(t, conn2, time.Second, testMessage))
	})

	t.Run("Client Count", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create test clients for different horses
		clients := []*Client{
			{hub: hub, send: make(chan []byte, 256), horseID: 1},
			{hub: hub, send: make(chan []byte, 256), horseID: 1},
			{hub: hub, send: make(chan []byte, 256), horseID: 2},
		}

		// Register clients
		for _, client := range clients {
			hub.register <- client
		}
		time.Sleep(100 * time.Millisecond)

		// Verify client counts
		assert.Equal(t, 2, hub.GetActiveClientsCount(1))
		assert.Equal(t, 1, hub.GetActiveClientsCount(2))
		assert.Equal(t, 0, hub.GetActiveClientsCount(3))
	})

	t.Run("Connection Timeout", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create test client with old ping time
		client := &Client{
			hub:      hub,
			send:     make(chan []byte, 256),
			horseID:  1,
			lastPing: time.Now().Add(-2 * pongWait),
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)

		// Verify client is removed after timeout
		hub.mu.RLock()
		clients := hub.clients[client.horseID]
		hub.mu.RUnlock()

		assert.Empty(t, clients)
	})

	t.Run("Message Buffer Overflow", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create client with small buffer
		client := &Client{
			hub:     hub,
			send:    make(chan []byte, 1),
			horseID: 1,
		}

		hub.register <- client
		time.Sleep(100 * time.Millisecond)

		// Send multiple messages to overflow buffer
		for i := 0; i < 10; i++ {
			hub.BroadcastToHorse(1, map[string]string{"msg": "overflow"})
		}
		time.Sleep(100 * time.Millisecond)

		// Verify client is removed due to buffer overflow
		hub.mu.RLock()
		clients := hub.clients[client.horseID]
		hub.mu.RUnlock()

		assert.Empty(t, clients)
	})
}

func TestWebSocketConnection(t *testing.T) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	t.Run("Connection Lifecycle", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		// Create test server
		server := testutil.CreateTestWebSocketServer(hub)
		defer server.Close()

		// Connect client
		conn, err := testutil.CreateTestWebSocketClient(server, "/ws/1")
		require.NoError(t, err)
		defer conn.Close()

		// Send ping
		err = conn.WriteMessage(websocket.PingMessage, nil)
		require.NoError(t, err)

		// Verify pong response
		_, _, err = conn.ReadMessage()
		require.NoError(t, err)
	})

	t.Run("Message Size Limit", func(t *testing.T) {
		hub := NewHub()
		go hub.Run()

		server := testutil.CreateTestWebSocketServer(hub)
		defer server.Close()

		conn, err := testutil.CreateTestWebSocketClient(server, "/ws/1")
		require.NoError(t, err)
		defer conn.Close()

		// Send oversized message
		bigMessage := make([]byte, maxMessageSize+1)
		err = conn.WriteMessage(websocket.TextMessage, bigMessage)
		assert.Error(t, err)
	})
}
