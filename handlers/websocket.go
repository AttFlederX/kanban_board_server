package handlers

import (
	"errors"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Client represents a websocket client connection
type Client struct {
	Conn   *websocket.Conn
	UserID primitive.ObjectID
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients mapped by user ID
	clients map[primitive.ObjectID]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

// Message represents a websocket message about task changes
type Message struct {
	Type   string      `json:"type"` // "create", "update", "delete"
	TaskID string      `json:"taskId"`
	UserID string      `json:"userId"`
	Data   interface{} `json:"data,omitempty"`
}

var hub *Hub

// JWTSecret stores the JWT secret for token validation
var JWTSecret string

// SetJWTSecret sets the JWT secret for websocket authentication
func SetJWTSecret(secret string) {
	JWTSecret = secret
}

// Claims represents JWT token claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	GoogleID string `json:"google_id"`
	jwt.RegisteredClaims
}

// validateTokenAndGetUserID validates JWT token and returns user ID
func validateTokenAndGetUserID(tokenString string) (primitive.ObjectID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return primitive.NilObjectID, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid token claims")
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user ID in token")
	}

	return userID, nil
}

// InitHub initializes the websocket hub
func InitHub() {
	hub = &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[primitive.ObjectID]map[*Client]bool),
	}
	go hub.run()
}

// run handles hub operations
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered for user %s", client.UserID.Hex())

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					client.Conn.Close()
					if len(clients) == 0 {
						delete(h.clients, client.UserID)
					}
					log.Printf("Client unregistered for user %s", client.UserID.Hex())
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Convert message UserID to ObjectID
			userObjectID, err := primitive.ObjectIDFromHex(message.UserID)
			if err != nil {
				log.Printf("Invalid user ID in broadcast message: %s", message.UserID)
				continue
			}

			h.mu.RLock()
			clients := h.clients[userObjectID]
			h.mu.RUnlock()

			// Send message to all clients of the user
			for client := range clients {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error writing to client: %v", err)
					h.unregister <- client
				}
			}
		}
	}
}

// BroadcastTaskChange broadcasts a task change to all connected clients of a user
func BroadcastTaskChange(messageType string, taskID primitive.ObjectID, userID primitive.ObjectID, data interface{}) {
	if hub == nil {
		log.Println("Warning: Hub not initialized, cannot broadcast message")
		return
	}

	message := Message{
		Type:   messageType,
		TaskID: taskID.Hex(),
		UserID: userID.Hex(),
		Data:   data,
	}

	hub.broadcast <- message
}

// HandleWebSocket handles websocket connections
func HandleWebSocket(c *websocket.Conn) {
	// Get token from query params for web clients (browsers can't send custom headers)
	token := c.Query("token")
	if token == "" {
		log.Println("WebSocket connection rejected: missing token")
		c.Close()
		return
	}

	// Validate token and extract user ID
	userID, err := validateTokenAndGetUserID(token)
	if err != nil {
		log.Printf("WebSocket connection rejected: %v", err)
		c.Close()
		return
	}

	client := &Client{
		Conn:   c,
		UserID: userID,
	}

	hub.register <- client

	// Keep connection alive and handle disconnection
	defer func() {
		hub.unregister <- client
	}()

	for {
		// Read messages from client (for keep-alive pings)
		_, _, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}
