package role

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	GetRolePermissionsList() (*models.RolePermissionsList, error)
	CreateRole(req *models.CreateRoleRequest) (*models.Role, error)
	GetRolePermissionsByID(id int) (*models.RolePermissions, error)
	UpdateRolePermissionsByID(req *models.UpdateRolePermissionsRequest) (*models.RolePermissions, error)
	DeleteRolePermissions(req *models.DeleteRolePermissionsRequest) error
}

type RoleController struct {
	RoleService IRoleService
}

func NewRoleController(roleService IRoleService) *RoleController {
	return &RoleController{
		RoleService: roleService,
	}
}

func (contl *RoleController) RegisterRoutes(rg *gin.RouterGroup) {
	roleGroup := rg.Group("/roles")
	roleGroup.GET("/", contl.GetRolePermissionsList)
	roleGroup.POST("/", contl.CreateRole)
	roleGroup.PUT("/", contl.UpdateRolePermissions)
	roleGroup.DELETE("/", contl.DeleteRolePermissions)
	roleGroup.GET("/:id", contl.GetRolePermissionsById)
}

// DeleteRolePermissions godoc
//
//	@Summary		Delete roles
//	@Description	Remove roles
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.DeleteRolePermissionsRequest	true	"Delete roles request"
//	@Success		200		{object}	models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Failure		500		{object}	models.APIResponse
//	@Router			/roles [delete]
func (contl *RoleController) DeleteRolePermissions(ctx *gin.Context) {
	var req models.DeleteRolePermissionsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	err := contl.RoleService.DeleteRolePermissions(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, nil)
}

// UpdateRolePermissions godoc
//
//	@Summary	Update a role
//	@Tags		Role
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.UpdateRolePermissionsRequest	true	"Role update request"
//	@Success	200		{object}	models.APIResponse{data=models.RolePermissions}
//	@Failure	400		{object}	models.APIResponse
//	@Failure	500		{object}	models.APIResponse
//	@Router		/roles [put]
func (contl *RoleController) UpdateRolePermissions(ctx *gin.Context) {
	var req models.UpdateRolePermissionsRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	rolePermissions, err := contl.RoleService.UpdateRolePermissionsByID(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, rolePermissions)
}

// GetRoleListById godoc
//
//	@Summary	Get role by ID
//	@Tags		Role
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Role ID"
//	@Success	200	{object}	models.APIResponse{data=models.RolePermissions}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/roles/{id} [get]
func (contl *RoleController) GetRolePermissionsById(ctx *gin.Context) {
	var req models.GetRolePermissionsByIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(lib.InvalidParamRequestError.WithStack(err.Error()))
		return
	}

	rolePermissions, err := contl.RoleService.GetRolePermissionsByID(req.ID)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, rolePermissions)
}

// GetRoleList godoc
//
//	@Summary	Get list of roles
//	@Tags		Role
//	@Produce	json
//	@Success	200	{object}	models.APIResponse{data=models.RolePermissionsList}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/roles [get]
func (contl *RoleController) GetRolePermissionsList(ctx *gin.Context) {
	roleList, err := contl.RoleService.GetRolePermissionsList()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	lib.SendSucceedResponse(ctx, roleList)
}

// CreateRole godoc
//
//	@Summary	Create a new role
//	@Tags		Role
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.CreateRoleRequest	true	"Role creation request"
//	@Success	200		{object}	models.APIResponse{data=models.Role}
//	@Failure	400		{object}	models.APIResponse
//	@Failure	500		{object}	models.APIResponse
//	@Router		/roles [post]
func (contl *RoleController) CreateRole(ctx *gin.Context) {
	var req models.CreateRoleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	role, err := contl.RoleService.CreateRole(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, role)

}
