package cache

import "github.com/redis/go-redis/v9"

type Service struct {
}

func NewCacheServices(redisClient *redis.Client) *Service {
	return &Service{}
}
