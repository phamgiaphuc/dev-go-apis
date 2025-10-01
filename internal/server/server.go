package server

import (
	"fmt"
	"net/http"
	"time"

	"dev-go-apis/docs"
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/router"

	"github.com/go-redis/redis_rate/v10"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func NewServer(dbClient *sqlx.DB, cacheClient *redis.Client, rateLimiter *redis_rate.Limiter) *http.Server {
	router := router.NewRouter(dbClient, cacheClient, rateLimiter)

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", lib.PORT)
	docs.SwaggerInfo.BasePath = "/api"

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", lib.PORT),
		Handler:      router.InitRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
