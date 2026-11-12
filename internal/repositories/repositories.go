package repositories

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repositories struct {
	Users    *UserRepository
	Sessions *SessionRepository
	Auth     *AuthRepository
}

func NewRepositories(db *gorm.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		Users:    NewUserRepository(db),
		Sessions: NewSessionRepository(db),
		Auth:     NewAuthRepository(redisClient),
	}
}
