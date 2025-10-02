package cache

import (
	"context"
	"encoding/json"
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

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return repo.CacheClient.Set(ctx, key, data, expiration).Err()
}

func (repo *CacheRepository) GetKeys(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var keys []string

	iter := repo.CacheClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (repo *CacheRepository) GetValue(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := repo.CacheClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(val, dest)
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
