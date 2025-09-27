package middleware

import (
	"dev-go-apis/internal/lib"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Split(ctx.Request.Header.Get("Authorization"), " ")[1]
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		userWithClaims, err := lib.ParseToken(token, lib.ACCESS_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		ctx.Set("user", userWithClaims)
		ctx.Next()
	}
}

func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Request.Cookie("rt")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		userWithClaims, err := lib.ParseToken(strings.Split(token.String(), "=")[1], lib.REFRESH_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		ctx.Set("user", userWithClaims)
		ctx.Next()
	}
}
