package lib

import (
	"github.com/gin-gonic/gin"
)

var (
	// Server
	GIN_MODE   = GetEnvString("GIN_MODE", gin.DebugMode)
	PORT       = GetEnvInt("PORT", 8000)
	SERVER_URL = GetEnvString("SERVER_URL", "http://localhost:8000")
	CLIENT_URL = GetEnvString("CLIENT_URL", "http://localhost:3000")
	API_KEY    = GetEnvString("API_KEY", "@secret123")

	// Database
	MIGRATION_MODE = GetEnvInt("MIGRATION_MODE", 0)
	DATABASE_URL   = GetEnvString("DATABASE_URL", "")
	REDIS_URL      = GetEnvString("REDIS_URL", "")

	// Secret
	ACCESS_TOKEN_SECRET  = GetEnvString("ACCESS_TOKEN_SECRET", "@secret123")
	REFRESH_TOKEN_SECRET = GetEnvString("REFRESH_TOKEN_SECRET", "@secret123")
	ACCESS_TOKEN_TTL     = GetEnvString("ACCESS_TOKEN_TTL", "15m")
	REFRESH_TOKEN_TTL    = GetEnvString("REFRESH_TOKEN_TTL", "7d")
)
