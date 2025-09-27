package database

import (
	"context"
	"dev-go-apis/internal/lib"
	"log"
	"time"

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
