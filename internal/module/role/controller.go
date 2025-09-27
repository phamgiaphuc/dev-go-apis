package role

import (
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	GetRoleList() ([]models.RoleList, error)
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
}

// GetRoleList godoc
//
//	@Summary	Get list of roles
//	@Tags		Role
//	@Produce	json
//	@Success	200	{object}	models.Response
//	@Router		/roles [get]
func (contl *RoleController) GetRoleList(ctx *gin.Context) {
	roleList, err := contl.RoleService.GetRoleList()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Data: roleList,
	})
}
