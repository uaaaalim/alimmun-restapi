package cache

import "github.com/redis/go-redis/v9"

type CacheService struct {
}

func NewCacheServices(redisClient *redis.Client) *CacheService {
	return &CacheService{}
}
