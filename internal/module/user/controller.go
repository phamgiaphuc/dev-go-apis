package user

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/middleware"
	"dev-go-apis/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error)
	GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error)
}

type ICacheService interface {
	SetValue(string, interface{}, time.Duration) error
}

type UserController struct {
	Service      IUserService
	CacheService ICacheService
}

type GetUserByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func NewUserController(service IUserService, cacheService ICacheService) *UserController {
	return &UserController{
		Service:      service,
		CacheService: cacheService,
	}
}

func (contl *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	userGroup.GET("/:id", middleware.ApiHmacHandler(), contl.GetUserById)
	userGroup.GET("/me",
		middleware.AccessTokenHandler(),
		middleware.PermissionHandler(
			contl.Service,
			contl.CacheService,
			lib.AdminDashboard,
		),
		contl.GetMe,
	)
}

// Get user info godoc
//
//	@Summary	Get user info from token
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.APIResponse{data=models.GetMeResponse}
//	@Router		/users/me [get]
//	@Security	Bearer
func (contl *UserController) GetMe(ctx *gin.Context) {
	userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)
	lib.SendSucceedResponse(ctx, &models.GetMeResponse{
		User:      userWithClaims.User,
		SessionID: userWithClaims.SessionID,
	})
}

// GetUserById godoc
//
//	@Summary	Get user by ID
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	models.APIResponse{data=models.UserWithAccounts}
//	@Router		/users/{id} [get]
func (contl *UserController) GetUserById(ctx *gin.Context) {
	var req GetUserByIDRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := contl.Service.GetUserByID(userID)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}
	lib.SendSucceedResponse(ctx, user)
}
