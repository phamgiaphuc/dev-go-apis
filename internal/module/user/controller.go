package user

import (
	"dev-go-apis/internal/middleware"
	"dev-go-apis/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IUserService interface {
	GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error)
}

type UserController struct {
	Service IUserService
}

type GetUserByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func NewUserController(service IUserService) *UserController {
	return &UserController{
		Service: service,
	}
}

func (contl *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	userGroup.GET("/:id", contl.GetUserById)
	userGroup.GET("/me", middleware.AccessTokenMiddleware(), contl.GetMe)
}

// Get user info godoc
//
//	@Summary	Get user info from token
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Router		/users/me [get]
//
//	@Security	Bearer
func (contl *UserController) GetMe(ctx *gin.Context) {
	userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)
	ctx.JSON(http.StatusOK, &models.UserResponse{
		User: &userWithClaims.User,
	})
}

// GetUserById godoc
//
//	@Summary	Get user by ID
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	models.UserWithAccounts
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
	ctx.JSON(http.StatusOK, user)

}
