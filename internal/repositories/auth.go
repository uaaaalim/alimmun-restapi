package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	redisClient *redis.Client
}

func NewAuthRepository(redisClient *redis.Client) *AuthRepository {
	return &AuthRepository{redisClient: redisClient}
}

func (repo *AuthRepository) IncrementAttempts(
	ctx context.Context,
	key string,
	window time.Duration,
) (int, time.Duration, error) {
	tx := repo.redisClient.TxPipeline()

	count := tx.Incr(ctx, key)
	tx.ExpireNX(ctx, key, window)
	ttl := tx.PTTL(ctx, key)

	if _, err := tx.Exec(ctx); err != nil {
		return 0, 0, fmt.Errorf("increment attempts: %w", err)
	}

	return int(count.Val()), ttl.Val(), nil
}
