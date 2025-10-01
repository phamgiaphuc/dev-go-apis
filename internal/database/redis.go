package database

import (
	"context"
	"log"
	"time"

	"dev-go-apis/internal/lib"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	redisOpts, err := redis.ParseURL(lib.REDIS_URL)
	if err != nil {
		log.Fatalln(err)
	}
	redisClient := redis.NewClient(redisOpts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("🎉 Redis is connected")

	return redisClient
}

func NewRateLimter(client *redis.Client) *redis_rate.Limiter {
	limiter := redis_rate.NewLimiter(client)
	return limiter
}
