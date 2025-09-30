package lib

import (
	"github.com/gin-gonic/gin"
)

var (
	// Server
	GIN_MODE        = GetEnvString("GIN_MODE", gin.DebugMode)
	PORT            = GetEnvInt("PORT", 8000)
	SERVER_URL      = GetEnvString("SERVER_URL", "http://localhost:8000")
	CLIENT_URL      = GetEnvString("CLIENT_URL", "http://localhost:3000")
	API_KEY         = GetEnvString("API_KEY", "@secret123")
	HMAC_SECRET_KEY = GetEnvString("HMAC_SECRET_KEY", "@secret123")

	// Cors
	CORS_ALLOWED_ORIGINS = GetEnvStrings("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"})

	// Google OAuth2
	GOOGLE_CLIENT_ID     = GetEnvString("GOOGLE_CLIENT_ID", "")
	GOOGLE_CLIENT_SECRET = GetEnvString("GOOGLE_CLIENT_SECRET", "")
	GOOGLE_REDIRECT_URL  = GetEnvString("GOOGLE_REDIRECT_URL", "http://localhost:8000/api/auth/google/callback")

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
