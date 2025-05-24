package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type     string   `json:"type"`
	Username string   `json:"username,omitempty"`
	Status   string   `json:"status,omitempty"`
	Users    []string `json:"users,omitempty"`
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) getOnlineUsers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]string, 0, len(m.clients))
	for client := range m.clients {
		users = append(users, client.Username)
	}
	return users
}

func (m *Manager) broadcastOnlineUsers() {
	users := m.getOnlineUsers()
	message := Message{
		Type:  "online_users",
		Users: users,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	m.broadcast <- data
}

func (m *Manager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	client := &Client{
		Conn: conn,
	}

	// Read the initial message to get the username
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("error reading message: %v", err)
		conn.Close()
		return
	}

	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("error unmarshaling message: %v", err)
		conn.Close()
		return
	}

	client.Username = msg.Username
	m.register <- client

	go m.readPump(client)
}

func (m *Manager) readPump(client *Client) {
	defer func() {
		m.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()
			m.broadcastUserStatus(client.Username, true)
			m.broadcastOnlineUsers()

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.Conn.Close()
			}
			m.mu.Unlock()
			m.broadcastUserStatus(client.Username, false)
			m.broadcastOnlineUsers()

		case message := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("error: %v", err)
					client.Conn.Close()
					delete(m.clients, client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *Manager) broadcastUserStatus(username string, online bool) {
	status := "online"
	if !online {
		status = "offline"
	}

	message := Message{
		Type:     "user_status",
		Username: username,
		Status:   status,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshaling message: %v", err)
		return
	}

	m.broadcast <- data
}
