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

func (repo *CacheRepository) GetValue(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return repo.CacheClient.Get(ctx, key).Result()
}

func (repo *CacheRepository) DeleteValue(keys ...string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.CacheClient.Del(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	return true, err
}
