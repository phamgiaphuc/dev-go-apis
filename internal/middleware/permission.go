package middleware

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error)
}

func PermissionHandler(userRepo IUserService, requiredPermissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(requiredPermissions) == 0 {
			ctx.AbortWithError(http.StatusForbidden, lib.ResourceForbiddenError)
			return
		}

		userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)

		userPermissions, err := userRepo.GetUserPermissionsByRoleID(userWithClaims.User.RoleID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, lib.InternalServerError)
			ctx.Abort()
			return
		}

		perms := make([]string, 0, len(userPermissions))
		for _, up := range userPermissions {
			perms = append(perms, up.Name)
		}

		hasPermission := slices.ContainsFunc(requiredPermissions, func(r string) bool {
			return slices.Contains(perms, r)
		})

		if !hasPermission {
			ctx.AbortWithError(http.StatusForbidden, lib.ResourceForbiddenError)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
