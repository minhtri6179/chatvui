package models

import (
	"time"
)

// User represents a user in the chat system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUser creates a new user with the given username
func NewUser(username string) *User {
	return &User{
		ID:        generateID(), // This would normally use a UUID library
		Username:  username,
		CreatedAt: time.Now(),
	}
}

// Simple ID generator for demo purposes
// In a real application, use a proper UUID library
func generateID() string {
	return time.Now().Format("20060102150405") + RandomString(8)
}

// RandomString generates a random string of the specified length
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(result)
}
