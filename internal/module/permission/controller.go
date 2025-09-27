package permission

import (
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IPermissionService interface {
	GetPermissionList() ([]models.PermissionList, error)
	CreatePermissionGroup(body *models.CreatePermissionGroupRequest) (*models.PermissionGroup, error)
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
	permissionGroup.POST("/group", contl.CreatePermissionGroup)
}

// GetPermissionList godoc
//
//	@Summary	Get list of permissions by group
//	@Tags		Permission
//	@Produce	json
//	@Success	200	{object}	models.Response
//	@Router		/permissions [get]
func (contl *PermissionController) GetPermissionList(ctx *gin.Context) {
	permissionList, err := contl.PermissionService.GetPermissionList()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Data: permissionList,
	})
}

// CreatePermissionGroup godoc
//
//	@Summary	Create a permission group
//	@Tags		Permission
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.CreatePermissionGroupRequest	true	"Create a permission group"
//	@Success	200		{object}	models.Response
//	@Router		/permissions/group [post]
func (contl *PermissionController) CreatePermissionGroup(ctx *gin.Context) {
	var req models.CreatePermissionGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	permissionGroup, err := contl.PermissionService.CreatePermissionGroup(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Data: permissionGroup,
	})
}
