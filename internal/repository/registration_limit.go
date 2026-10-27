package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
)

type RegistrationLimitRepository struct {
	client *redis.Client
	prefix string
}

func NewRegistrationLimitRepository(connection *database.RedisConnection) *RegistrationLimitRepository {
	return &RegistrationLimitRepository{client: connection.Redis, prefix: connection.Prefix}
}

func (repo *RegistrationLimitRepository) key(ip string) string {
	return fmt.Sprintf("%s:auth:register:bans:%x", repo.prefix, sha256.Sum256([]byte(ip)))
}

func (repo *RegistrationLimitRepository) BanRemaining(ctx context.Context, ip string) (time.Duration, error) {
	remaining, err := repo.client.PTTL(ctx, repo.key(ip)).Result()
	if err != nil {
		return 0, fmt.Errorf("get registration ban ttl: %w", err)
	}
	switch remaining {
	case -2:
		return 0, nil
	case -1:
		return 0, fmt.Errorf("registration ban has no expiration")
	case 0:
		return time.Millisecond, nil
	}
	return remaining, nil
}

func (repo *RegistrationLimitRepository) Ban(ctx context.Context, ip string, duration time.Duration) error {
	if err := repo.client.Set(ctx, repo.key(ip), "1", duration).Err(); err != nil {
		return fmt.Errorf("set registration ban: %w", err)
	}
	return nil
}
