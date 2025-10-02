package lib

import (
	"net/http"

	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
)

func SendSucceedResponse(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &models.APIResponse{
		Success: true,
		Data:    data,
	})
}

func SendErrorResponse(ctx *gin.Context, err *models.APIError) {
	resp := &models.APIResponse{
		Success: false,
		Message: err.Message,
		Errors:  err.Errors,
	}
	if GIN_MODE != gin.ReleaseMode {
		resp.Stack = err.Stack
	}
	ctx.AbortWithStatusJSON(err.Code, resp)
}
