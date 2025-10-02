package role

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	GetRoleList() (*models.RoleList, error)
	CreateRole(req *models.CreateRoleRequest) (*models.RoleWithPermissions, error)
	UpdateRoleById(req *models.UpdateRoleRequest) (*models.RoleWithPermissions, error)
	DeleteRole(req *models.DeleteRolesRequest) error
	GetRoleById(req *models.GetRoleByIdRequest) (*models.RoleWithPermissions, error)
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
	roleGroup.GET("", contl.GetRoleList)
	roleGroup.POST("", contl.CreateRole)
	roleGroup.PUT("", contl.UpdateRole)
	roleGroup.DELETE("", contl.DeleteRole)
	roleGroup.GET("/:id", contl.GetRoleById)
}

// GetRoleList godoc
//
//	@Summary	Get list of roles
//	@Tags		Role & Permission
//	@Produce	json
//	@Success	200	{object}	models.APIResponse{data=models.RoleList}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/roles [get]
func (contl *RoleController) GetRoleList(ctx *gin.Context) {
	roleList, err := contl.RoleService.GetRoleList()
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, roleList)
}

// CreateRole godoc
//
//	@Summary	Create a new role
//	@Tags		Role & Permission
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.CreateRoleRequest	true	"Role creation request"
//	@Success	200		{object}	models.APIResponse{data=models.RoleWithPermissions}
//	@Failure	400		{object}	models.APIResponse
//	@Failure	500		{object}	models.APIResponse
//	@Router		/roles [post]
func (contl *RoleController) CreateRole(ctx *gin.Context) {
	var req models.CreateRoleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	roleWithPermissions, err := contl.RoleService.CreateRole(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, roleWithPermissions)
}

// UpdateRole godoc
//
//	@Summary	Update a role
//	@Tags		Role & Permission
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.UpdateRoleRequest	true	"Role update request"
//	@Success	200		{object}	models.APIResponse{data=models.RoleWithPermissions}
//	@Failure	400		{object}	models.APIResponse
//	@Failure	500		{object}	models.APIResponse
//	@Router		/roles [put]
func (contl *RoleController) UpdateRole(ctx *gin.Context) {
	var req models.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	roleWithPermissions, err := contl.RoleService.UpdateRoleById(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, roleWithPermissions)
}

// DeleteRole godoc
//
//	@Summary		Delete roles
//	@Description	Remove roles
//	@Tags			Role & Permission
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.DeleteRolesRequest	true	"Delete roles request"
//	@Success		200		{object}	models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Failure		500		{object}	models.APIResponse
//	@Router			/roles [delete]
func (contl *RoleController) DeleteRole(ctx *gin.Context) {
	var req models.DeleteRolesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	err := contl.RoleService.DeleteRole(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, nil)
}

// GetRoleById godoc
//
//	@Summary	Get role by ID
//	@Tags		Role & Permission
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Role ID"
//	@Success	200	{object}	models.APIResponse{data=models.RoleWithPermissions}
//	@Failure	400	{object}	models.APIResponse
//	@Failure	500	{object}	models.APIResponse
//	@Router		/roles/{id} [get]
func (contl *RoleController) GetRoleById(ctx *gin.Context) {
	var req models.GetRoleByIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(lib.InvalidParamRequestError.WithStack(err.Error()))
		return
	}

	roleWithPermissions, err := contl.RoleService.GetRoleById(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, roleWithPermissions)
}
