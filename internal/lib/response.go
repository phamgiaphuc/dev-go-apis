package lib

import (
	"dev-go-apis/internal/models"
	"net/http"

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
	}
	if GIN_MODE != gin.ReleaseMode {
		resp.Stack = err.GetStack()
	}
	ctx.JSON(err.Code, resp)
}
