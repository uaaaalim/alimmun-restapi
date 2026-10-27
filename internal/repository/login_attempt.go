package repository

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uaaaalim/alimmun-restapi/internal/database"
)

type LoginAttemptRepository struct {
	Redis       *redis.Client
	RedisPrefix string
}

func NewLoginAttemptRepository(redisConnection *database.RedisConnection) *LoginAttemptRepository {
	return &LoginAttemptRepository{Redis: redisConnection.Redis, RedisPrefix: redisConnection.Prefix}
}

func (repo *LoginAttemptRepository) keys(ip, email, action string) string {
	id := fmt.Sprintf("%x", sha256.Sum256(
		[]byte(ip+"\x00"+email),
	))

	return repo.RedisPrefix + ":auth:" + action + ":" + id
}

func (repo *LoginAttemptRepository) BanRemaining(
	ctx context.Context,
	ip, email string,
) (time.Duration, error) {
	key := repo.keys(ip, email, "login:bans")

	remaining, err := repo.Redis.PTTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("get login ban ttl: %w", err)
	}

	return remaining, nil
}

func (repo *LoginAttemptRepository) ResetFailures(ctx context.Context, ip, email string) error {
	key := repo.keys(ip, email, "login:attempts")

	if err := repo.Redis.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *LoginAttemptRepository) IncrementLoginAttempt(ctx context.Context, ip, email string) (int64, error) {
	count, err := repo.Redis.Incr(ctx, repo.keys(ip, email, "login:attempts")).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *LoginAttemptRepository) ExpireAttempts(
	ctx context.Context,
	ip, email string,
	duration time.Duration,
) error {
	if err := repo.Redis.Expire(ctx, repo.keys(ip, email, "login:attempts"), duration).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *LoginAttemptRepository) BanAttempts(
	ctx context.Context, ip, email string, duration time.Duration,
) error {
	if err := repo.Redis.Set(ctx, repo.keys(ip, email, "login:bans"), "1", duration).Err(); err != nil {
		return err
	}

	return nil
}
