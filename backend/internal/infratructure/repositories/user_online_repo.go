package repositories

import (
	"backend/internal/domain/user"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOnlineRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisOnlineRepo(client *redis.Client, ctx context.Context) *RedisOnlineRepo {
	return &RedisOnlineRepo{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisOnlineRepo) SetUserOnline(userID, username string) error {
	status := user.OnlineStatus{
		UserID:   userID,
		Username: username,
		LastSeen: time.Now(),
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return err
	}

	// Store user status with expiration
	return r.client.Set(r.ctx, userID, string(statusJSON), 10*time.Second).Err()
}
