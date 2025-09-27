package middleware

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error)
}

type ICacheService interface {
	SetValue(string, interface{}, time.Duration) error
}

func PermissionHandler(userRepo IUserService, cacheSerivce ICacheService, requiredPermissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(requiredPermissions) == 0 {
			ctx.AbortWithError(http.StatusForbidden, lib.ResourceForbiddenError)
			return
		}

		userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)

		userPermissions, err := userRepo.GetUserPermissionsByRoleID(userWithClaims.User.RoleID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, lib.InternalServerError)
			return
		}

		perms := make([]string, 0, len(userPermissions))
		for _, up := range userPermissions {
			perms = append(perms, up.Name)
		}

		values, _ := json.Marshal(perms)

		err = cacheSerivce.SetValue(fmt.Sprintf("role:%s", userWithClaims.User.Role), values, 0)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, lib.InternalServerError)
			return
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
