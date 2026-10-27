package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	Redis  *redis.Client
	Prefix string
}

func ConnectRedis(redisUrl, redisPrefix string) (*RedisConnection, error) {
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url")
	}

	options.DialTimeout = 5 * time.Second
	options.ReadTimeout = 3 * time.Second
	options.WriteTimeout = 3 * time.Second

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisConnection{Redis: client, Prefix: redisPrefix}, nil
}
