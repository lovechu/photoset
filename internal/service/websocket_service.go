package service

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket message types
const (
	WSTypeMessage    = "message"    // New private message
	WSTypeNotification = "notification" // New notification
	WSTypeUnreadCount  = "unread_count" // Unread count update
	WSTypeOnlineStatus = "online_status" // User online/offline status
	WSTypeTyping       = "typing"       // Typing indicator
	WSTypePing         = "ping"         // Heartbeat ping
	WSTypePong         = "pong"         // Heartbeat pong
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WSClient represents a single WebSocket connection
type WSClient struct {
	hub    *Hub
	conn   *websocket.Conn
	userID uint
	send   chan []byte
	mu     sync.Mutex
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients indexed by user ID
	clients map[uint]map[*WSClient]bool

	// Register requests from clients
	register chan *WSClient

	// Unregister requests from clients
	unregister chan *WSClient

	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*WSClient]bool),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*WSClient]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()

			// Notify others that this user is online
			h.broadcastOnlineStatus(client.userID, true)
			log.Printf("[WebSocket] User %d connected (total connections: %d)", client.userID, h.TotalConnections())

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
						// Notify others that this user is offline
						h.broadcastOnlineStatus(client.userID, false)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("[WebSocket] User %d disconnected", client.userID)
		}
	}
}

// SendToUser sends a message to a specific user (all their connections)
func (h *Hub) SendToUser(userID uint, msgType string, payload interface{}) {
	msg := WSMessage{
		Type:    msgType,
		Payload: payload,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WebSocket] Failed to marshal message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.send <- data:
			default:
				// Client send buffer full, skip
				log.Printf("[WebSocket] User %d send buffer full, skipping", userID)
			}
		}
	}
}

// IsUserOnline checks if a user has any active connections
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	return ok && len(clients) > 0
}

// TotalConnections returns total number of active connections
func (h *Hub) TotalConnections() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, clients := range h.clients {
		count += len(clients)
	}
	return count
}

// OnlineUserIDs returns a list of currently online user IDs
func (h *Hub) OnlineUserIDs() []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]uint, 0, len(h.clients))
	for userID := range h.clients {
		ids = append(ids, userID)
	}
	return ids
}

// broadcastOnlineStatus broadcasts a user's online status to all connected users
func (h *Hub) broadcastOnlineStatus(userID uint, online bool) {
	// For now, we don't broadcast to all users (could be expensive)
	// This can be enabled later for "who's online" features
	_ = userID
	_ = online
}

// NewClient creates a new WebSocket client
func (h *Hub) NewClient(conn *websocket.Conn, userID uint) *WSClient {
	return &WSClient{
		hub:    h,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
	}
}

// Register registers a client with the hub
func (h *Hub) Register(client *WSClient) {
	h.register <- client
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *WSClient) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WebSocket] Read error for user %d: %v", c.userID, err)
			}
			break
		}

		// Handle incoming messages (e.g., ping/pong, typing indicators)
		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err == nil {
			c.handleMessage(wsMsg)
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *WSClient) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Drain queued messages into current write
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming WebSocket messages
func (c *WSClient) handleMessage(msg WSMessage) {
	switch msg.Type {
	case WSTypePing:
		// Respond with pong
		pong := WSMessage{Type: WSTypePong, Payload: map[string]string{"status": "ok"}}
		data, _ := json.Marshal(pong)
		select {
		case c.send <- data:
		default:
		}
	case WSTypeTyping:
		// Forward typing indicator to the target user
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if targetID, ok := payload["to_user_id"].(float64); ok {
				c.hub.SendToUser(uint(targetID), WSTypeTyping, map[string]interface{}{
					"from_user_id": c.userID,
				})
			}
		}
	}
}
