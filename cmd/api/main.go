package main

import (
	"context"
	"dev-go-apis/internal/cache"
	"dev-go-apis/internal/database"
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/server"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

//	@title		Dev Go APIs
//	@version	v1.0.0
//	@host		localhost:8000
//	@BasePath	/api

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	if lib.MIGRATION_MODE == 1 {
		database.MigrateDB()
	}

	dbClient := database.NewDatabaseClient()
	defer dbClient.Close()

	cacheClient := cache.NewRedisClient()
	defer cacheClient.Close()

	server := server.NewServer(dbClient, cacheClient)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	go func() {
		log.Printf("🚀 Server is running on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	<-done
}
