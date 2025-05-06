package user

type OnlineUserRepository interface {
	SetUserOnline(userID, username string) error
	SetUserOffline(userID string) error
	IsUserOnline(userID string) (bool, error)
	GetOnlineUsers() ([]OnlineStatus, error)
	CountOnlineUsers() (int, error)
	UpdateUserActivity(userID string) error
}
