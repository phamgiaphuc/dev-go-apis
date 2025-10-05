package middleware

import (
	"strings"

	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
)

func AccessTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenKey := ctx.Request.Header.Get("Authorization")
		if tokenKey == "" {
			lib.SendErrorResponse(ctx, lib.UnauthorizedError)
			return
		}
		token := strings.Split(ctx.Request.Header.Get("Authorization"), " ")[1]
		payload, err := lib.ParseToken[models.JwtUserPayload](token, lib.ACCESS_TOKEN_SECRET)
		if err != nil {
			lib.SendErrorResponse(ctx, lib.UnauthorizedError)
			return
		}
		ctx.Set("user", payload.UserWithClaims)
		ctx.Next()
	}
}

func RefreshTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Request.Cookie("rt")
		if err != nil {
			lib.SendErrorResponse(ctx, lib.UnauthorizedError)
			return
		}
		payload, err := lib.ParseToken[models.JwtUserPayload](strings.Split(token.String(), "=")[1], lib.REFRESH_TOKEN_SECRET)
		if err != nil {
			lib.SendErrorResponse(ctx, lib.UnauthorizedError)
			return
		}
		ctx.Set("user", payload.UserWithClaims)
		ctx.Next()
	}
}
