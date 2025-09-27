package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	CacheClient *redis.Client
}

func NewCacheRepository(cacheClient *redis.Client) *CacheRepository {
	return &CacheRepository{
		CacheClient: cacheClient,
	}
}

func (repo *CacheRepository) SetValue(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return repo.CacheClient.Set(ctx, key, value, expiration).Err()
}
