package handlers

import (
	"net/http"
	"sync"

	"chatvui/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Handler contains all the dependencies for the API handlers
type Handler struct {
	hub *models.Hub
	// You would typically have a database here as well

	// In-memory user store for this example
	users     map[string]*models.User
	userMutex sync.RWMutex
}

// NewHandler creates a new handler with the given hub
func NewHandler(hub *models.Hub) *Handler {
	return &Handler{
		hub:   hub,
		users: make(map[string]*models.User),
	}
}

// RegisterUser registers a new user
func (h *Handler) RegisterUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Create a new user
	user := models.NewUser(request.Username)

	// Save user
	h.userMutex.Lock()
	h.users[user.ID] = user
	h.userMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetUsers returns a list of all users
func (h *Handler) GetUsers(c *gin.Context) {
	h.userMutex.RLock()
	userList := make([]*models.User, 0, len(h.users))
	for _, user := range h.users {
		userList = append(userList, user)
	}
	h.userMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
	})
}

// GetActiveUsers returns a list of all active users in the chat
func (h *Handler) GetActiveUsers(c *gin.Context) {
	activeUsers := h.hub.GetActiveUsers()

	c.JSON(http.StatusOK, gin.H{
		"users": activeUsers,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now (not secure for production)
	},
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get user
	h.userMutex.RLock()
	user, exists := h.users[userID]
	h.userMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Create a new client
	client := models.NewClient(h.hub, conn, user)

	// Register client
	h.hub.Register <- client

	// Start handling messages
	go client.WritePump()
	go client.ReadPump()
}
