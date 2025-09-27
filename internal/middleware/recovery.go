package middleware

import (
	"dev-go-apis/internal/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, lib.InternalServerError)
	})
}
