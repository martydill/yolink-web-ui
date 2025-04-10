package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go-svelte-app/backend/mqtt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type Client struct {
	conn *websocket.Conn
	send chan interface{}
}

var (
	clients = make(map[*Client]bool)
	mu      sync.RWMutex
)

// Handle WebSocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan interface{}, 100),
	}

	// Register client
	mu.Lock()
	clients[client] = true
	mu.Unlock()

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// Start the WebSocket update broadcaster
func StartBroadcaster() {
	updateChannel := mqtt.GetUpdateChannel()
	for update := range updateChannel {
		broadcastUpdate(update)
	}
}

// Broadcast update to all connected clients
func broadcastUpdate(update mqtt.DeviceUpdate) {
	mu.RLock()
	defer mu.RUnlock()

	for client := range clients {
		select {
		case client.send <- update:
		default:
			// If client's send channel is full, close the connection
			close(client.send)
			delete(clients, client)
		}
	}
}

// Write messages to the WebSocket connection
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}
		}
	}
}

// Read messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading from WebSocket: %v", err)
			}
			break
		}

		// Handle incoming messages if needed
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing WebSocket message: %v", err)
			continue
		}
	}
} 