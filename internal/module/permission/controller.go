package permission

import (
	"net/http"

	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
)

type IPermissionService interface {
	GetPermissionList() ([]models.PermissionList, error)
}

type PermissionController struct {
	PermissionService IPermissionService
}

func NewPermissionController(permissionService IPermissionService) *PermissionController {
	return &PermissionController{
		PermissionService: permissionService,
	}
}

func (contl *PermissionController) RegisterRoutes(rg *gin.RouterGroup) {
	permissionGroup := rg.Group("/permissions")
	permissionGroup.GET("/", contl.GetPermissionList)
}

// GetPermissionList godoc
//
//	@Summary	Get list of permissions by group
//	@Tags		Role & Permission
//	@Produce	json
//	@Success	200	{object}	models.APIResponse
//	@Router		/permissions [get]
func (contl *PermissionController) GetPermissionList(ctx *gin.Context) {
	permissionList, err := contl.PermissionService.GetPermissionList()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.APIResponse{
		Data: permissionList,
	})
}
