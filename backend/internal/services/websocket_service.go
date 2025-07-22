package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// WebSocketService handles real-time notifications
type WebSocketService struct {
	clients    map[string]*Client
	broadcast  chan Notification
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
}

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	UserID   string
	Socket   *websocket.Conn
	Send     chan []byte
	Service  *WebSocketService
}

// Notification represents a notification message
type Notification struct {
	Type      string      `json:"type"`
	UserID    string      `json:"user_id,omitempty"`
	Title     string      `json:"title"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NotificationTypes
const (
	NotificationTypeTransaction = "transaction"
	NotificationTypeSecurity   = "security"
	NotificationTypeSystem     = "system"
	NotificationTypeKYC        = "kyc"
	NotificationTypeFloat      = "float"
)

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Notification),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

// Start starts the WebSocket service
func (ws *WebSocketService) Start() {
	for {
		select {
		case client := <-ws.register:
			ws.mutex.Lock()
			ws.clients[client.ID] = client
			ws.mutex.Unlock()
			log.Info().Str("client_id", client.ID).Str("user_id", client.UserID).Msg("Client connected")

		case client := <-ws.unregister:
			ws.mutex.Lock()
			if _, ok := ws.clients[client.ID]; ok {
				delete(ws.clients, client.ID)
				close(client.Send)
			}
			ws.mutex.Unlock()
			log.Info().Str("client_id", client.ID).Str("user_id", client.UserID).Msg("Client disconnected")

		case notification := <-ws.broadcast:
			ws.broadcastNotification(notification)
		}
	}
}

// HandleWebSocket handles WebSocket connections
func (ws *WebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get user ID from query parameters or headers
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = r.Header.Get("X-User-ID")
	}

	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection to WebSocket")
		return
	}

	// Create new client
	client := &Client{
		ID:      generateClientID(),
		UserID:  userID,
		Socket:  conn,
		Send:    make(chan []byte, 256),
		Service: ws,
	}

	// Register client
	ws.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// SendNotification sends a notification to a specific user
func (ws *WebSocketService) SendNotification(userID string, notification Notification) {
	notification.UserID = userID
	notification.Timestamp = time.Now()

	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	// Send to all clients of the user
	for _, client := range ws.clients {
		if client.UserID == userID {
			select {
			case client.Send <- ws.notificationToJSON(notification):
			default:
				close(client.Send)
				delete(ws.clients, client.ID)
			}
		}
	}

	log.Info().
		Str("user_id", userID).
		Str("type", notification.Type).
		Str("title", notification.Title).
		Msg("Notification sent to user")
}

// BroadcastNotification sends a notification to all connected clients
func (ws *WebSocketService) BroadcastNotification(notification Notification) {
	notification.Timestamp = time.Now()
	ws.broadcast <- notification
}

// SendTransactionNotification sends a transaction notification
func (ws *WebSocketService) SendTransactionNotification(userID, transactionType, amount, transactionID string) {
	notification := Notification{
		Type:    NotificationTypeTransaction,
		Title:   "Transaction Completed",
		Message: fmt.Sprintf("Your %s transaction of KES %s has been completed successfully.", transactionType, amount),
		Data: map[string]interface{}{
			"transaction_id": transactionID,
			"type":          transactionType,
			"amount":        amount,
		},
	}

	ws.SendNotification(userID, notification)
}

// SendSecurityNotification sends a security notification
func (ws *WebSocketService) SendSecurityNotification(userID, alertType string) {
	notification := Notification{
		Type:    NotificationTypeSecurity,
		Title:   "Security Alert",
		Message: fmt.Sprintf("Security alert: %s. If this wasn't you, please contact support immediately.", alertType),
		Data: map[string]interface{}{
			"alert_type": alertType,
		},
	}

	ws.SendNotification(userID, notification)
}

// SendKYCNotification sends a KYC notification
func (ws *WebSocketService) SendKYCNotification(userID, status, message string) {
	notification := Notification{
		Type:    NotificationTypeKYC,
		Title:   "KYC Status Update",
		Message: message,
		Data: map[string]interface{}{
			"kyc_status": status,
		},
	}

	ws.SendNotification(userID, notification)
}

// SendFloatNotification sends a float notification to agents
func (ws *WebSocketService) SendFloatNotification(agentID, message string) {
	notification := Notification{
		Type:    NotificationTypeFloat,
		Title:   "Float Update",
		Message: message,
		Data: map[string]interface{}{
			"agent_id": agentID,
		},
	}

	ws.SendNotification(agentID, notification)
}

// GetConnectedUsers returns the number of connected users
func (ws *WebSocketService) GetConnectedUsers() int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	return len(ws.clients)
}

// GetUserConnections returns the number of connections for a specific user
func (ws *WebSocketService) GetUserConnections(userID string) int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	count := 0
	for _, client := range ws.clients {
		if client.UserID == userID {
			count++
		}
	}
	return count
}

// broadcastNotification broadcasts a notification to all connected clients
func (ws *WebSocketService) broadcastNotification(notification Notification) {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	notificationJSON := ws.notificationToJSON(notification)

	for _, client := range ws.clients {
		select {
		case client.Send <- notificationJSON:
		default:
			close(client.Send)
			delete(ws.clients, client.ID)
		}
	}

	log.Info().
		Str("type", notification.Type).
		Str("title", notification.Title).
		Int("recipients", len(ws.clients)).
		Msg("Notification broadcasted")
}

// notificationToJSON converts a notification to JSON
func (ws *WebSocketService) notificationToJSON(notification Notification) []byte {
	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal notification to JSON")
		return []byte("{}")
	}
	return jsonData
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.Service.unregister <- c
		c.Socket.Close()
	}()

	c.Socket.SetReadLimit(512)
	c.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Socket.SetPongHandler(func(string) error {
		c.Socket.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Str("client_id", c.ID).Msg("WebSocket read error")
			}
			break
		}

		// Handle incoming messages (ping, etc.)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msgType, ok := msg["type"].(string); ok {
				switch msgType {
				case "ping":
					response := map[string]interface{}{
						"type": "pong",
						"timestamp": time.Now(),
					}
					responseJSON, _ := json.Marshal(response)
					c.Send <- responseJSON
				}
			}
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
} 