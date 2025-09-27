package server

import (
	"dev-go-apis/docs"
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/router"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func NewServer(dbClient *sqlx.DB, cacheClient *redis.Client) *http.Server {
	router := router.NewRouter(dbClient, cacheClient)

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
