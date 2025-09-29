package middleware

import (
	"dev-go-apis/internal/lib"

	"github.com/gin-gonic/gin"
)

func ApiKeyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.Request.Header.Get("api_key")
		if apiKey == "" || apiKey != lib.API_KEY {
			lib.SendErrorResponse(ctx, lib.MissingAPIKeyError)
			return
		}
		ctx.Next()
	}
}
