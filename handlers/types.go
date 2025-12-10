package handlers

import (
	"sync"

	"github.com/AttFlederX/kanban_board_server/models"
	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GoogleSignInRequest represents the request body for Google sign-in
type GoogleSignInRequest struct {
	IDToken string `json:"id_token"`
}

// AuthResponse represents the response for authentication
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

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

// Claims represents JWT token claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	GoogleID string `json:"google_id"`
	jwt.RegisteredClaims
}
