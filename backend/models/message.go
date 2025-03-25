package models

import (
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "user" or "system"
}

// NewMessage creates a new message
func NewMessage(text, userID, username, messageType string) *Message {
	return &Message{
		ID:        generateID(),
		Text:      text,
		UserID:    userID,
		Username:  username,
		Timestamp: time.Now(),
		Type:      messageType,
	}
}

// SystemMessage creates a new system message
func SystemMessage(text string) *Message {
	return &Message{
		ID:        generateID(),
		Text:      text,
		UserID:    "system",
		Username:  "System",
		Timestamp: time.Now(),
		Type:      "system",
	}
}
