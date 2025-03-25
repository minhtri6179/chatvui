package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Inbound messages from the clients
	Broadcast chan *Message

	// Mutex for concurrent access to clients map
	mu sync.Mutex

	// Messages history (last 50 messages)
	history []*Message
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		history:    make([]*Message, 0, 50),
	}
}

// Client represents a WebSocket connection and a user
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan *Message
	User     *User
	isActive bool
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// Send welcome message to the new client
			welcomeMsg := SystemMessage("Welcome to the chat, " + client.User.Username + "!")
			client.send <- welcomeMsg

			// Send join notification to everyone
			joinMsg := SystemMessage(client.User.Username + " has joined the chat.")
			h.Broadcast <- joinMsg

			// Send history messages to new client
			for _, msg := range h.history {
				client.send <- msg
			}

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Only broadcast leave message if client was previously active
				if client.isActive {
					// Send leave notification
					leaveMsg := SystemMessage(client.User.Username + " has left the chat.")
					h.Broadcast <- leaveMsg
				}
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			// Add message to history, keeping only the last 50
			h.addToHistory(message)

			// Broadcast to all clients
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// addToHistory adds a message to the history, keeping only the last 50
func (h *Hub) AddToHistory(message *Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.addToHistory(message)
}

// internal method (without locking)
func (h *Hub) addToHistory(message *Message) {
	// Only store regular messages, not system notifications
	if message.Type == "user" {
		if len(h.history) >= 50 {
			// Shift elements to remove the oldest
			h.history = append(h.history[1:], message)
		} else {
			h.history = append(h.history, message)
		}
	}
}

// GetActiveUsers returns a list of active users
func (h *Hub) GetActiveUsers() []*User {
	h.mu.Lock()
	defer h.mu.Unlock()

	users := make([]*User, 0, len(h.clients))
	for client := range h.clients {
		if client.isActive {
			users = append(users, client.User)
		}
	}
	return users
}

// SendMessage sends a new message and broadcasts it to all clients
func (h *Hub) SendMessage(message *Message) {
	h.Broadcast <- message
}
