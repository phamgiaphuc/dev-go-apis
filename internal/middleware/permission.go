package middleware

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	CheckUserRole(uuid.UUID, []string) (bool, error)
}

func PermissionHandler(userService IUserService, requiredPermissions []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(requiredPermissions) == 0 {
			ctx.Next()
			return
		}

		userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)
		hasPermission, err := userService.CheckUserRole(userWithClaims.UserID, requiredPermissions)
		if err != nil {
			lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
			return
		}
		if !hasPermission {
			lib.SendErrorResponse(ctx, lib.ResourceForbiddenError)
			return
		}

		ctx.Next()
	}
}
