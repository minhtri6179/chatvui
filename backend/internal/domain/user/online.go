package user

import "time"

type OnlineStatus struct {
	UserID   string
	Username string
	LastSeen time.Time
}
