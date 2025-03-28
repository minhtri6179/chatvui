package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

type UserStatus struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	LastSeen time.Time `json:"last_seen"`
	IsOnline bool      `json:"is_online"`
}

func NewRedisService(addr string) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisService{client: client}, nil
}

// UpdateUserStatus updates or creates a user's online status
func (s *RedisService) UpdateUserStatus(ctx context.Context, user UserStatus) error {
	key := fmt.Sprintf("user:%s", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetUserStatus retrieves a user's status
func (s *RedisService) GetUserStatus(ctx context.Context, userID string) (*UserStatus, error) {
	key := fmt.Sprintf("user:%s", userID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var status UserStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// GetOnlineUsers retrieves all online users
func (s *RedisService) GetOnlineUsers(ctx context.Context) ([]UserStatus, error) {
	pattern := "user:*"
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var users []UserStatus
	for _, key := range keys {
		data, err := s.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var status UserStatus
		if err := json.Unmarshal(data, &status); err != nil {
			continue
		}

		if status.IsOnline {
			users = append(users, status)
		}
	}

	return users, nil
}

// SetUserOffline marks a user as offline
func (s *RedisService) SetUserOffline(ctx context.Context, userID string) error {
	status, err := s.GetUserStatus(ctx, userID)
	if err != nil {
		return err
	}

	status.IsOnline = false
	status.LastSeen = time.Now()
	return s.UpdateUserStatus(ctx, *status)
}

// Close closes the Redis connection
func (s *RedisService) Close() error {
	return s.client.Close()
}
