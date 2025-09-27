package middleware

import (
	"dev-go-apis/internal/lib"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AccessTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := strings.Split(ctx.Request.Header.Get("Authorization"), " ")[1]
		if token == "" {
			ctx.AbortWithError(http.StatusUnauthorized, lib.UnauthorizedError)
			return
		}
		userWithClaims, err := lib.ParseToken(token, lib.ACCESS_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, lib.UnauthorizedError)
			return
		}
		ctx.Set("user", userWithClaims)
		ctx.Next()
	}
}

func RefreshTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Request.Cookie("rt")
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, lib.UnauthorizedError)
			return
		}
		userWithClaims, err := lib.ParseToken(strings.Split(token.String(), "=")[1], lib.REFRESH_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, lib.UnauthorizedError)
			return
		}
		ctx.Set("user", userWithClaims)
		ctx.Next()
	}
}
