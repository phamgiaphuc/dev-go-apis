package role

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	GetRoleList() ([]models.RoleList, error)
	CreateRole(req *models.CreateRoleRequest) (*models.Role, error)
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
	roleGroup.GET("/", contl.GetRoleList)
	roleGroup.POST("/", contl.CreateRole)
}

// GetRoleList godoc
//
//	@Summary	Get list of roles
//	@Tags		Role
//	@Produce	json
//	@Success	200	{object}	models.APIResponse
//	@Router		/roles [get]
func (contl *RoleController) GetRoleList(ctx *gin.Context) {
	roleList, err := contl.RoleService.GetRoleList()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	lib.SendSucceedResponse(ctx, &models.GetRoleListResponse{
		Roles: roleList,
	})
}

// CreateRole godoc
//
//	@Summary	Create a new role
//	@Tags		Role
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.CreateRoleRequest	true	"Role creation request"
//	@Success	200		{object}	models.APIResponse{data=models.CreateRoleResponse}
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

	lib.SendSucceedResponse(ctx, &models.CreateRoleResponse{
		Role: role,
	})

}
