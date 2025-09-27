package middleware

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
			switch e := err.(type) {
			case *models.APIError:
				lib.SendErrorResponse(ctx, e)
			default:
				ctx.JSON(http.StatusInternalServerError, lib.InternalServerError)
			}
			ctx.Abort()
		}
	}
}
