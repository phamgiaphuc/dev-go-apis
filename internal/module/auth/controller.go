package auth

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.User, error)
	GenerateJwtTokens(userWithClaims *models.UserWithClaims) (*models.JwtTokens, error)
}

type ISessionService interface {
	CreateSession(session *models.Session) (*models.Session, error)
}

type AuthController struct {
	AuthService    IAuthService
	SessionService ISessionService
}

func NewAuthController(authService IAuthService, sessionService ISessionService) *AuthController {
	return &AuthController{
		AuthService:    authService,
		SessionService: sessionService,
	}
}

func (contl *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/register", contl.Register)
	authGroup.POST("/login", contl.Login)
	authGroup.POST("/logout", contl.Logout)
}

// Register godoc
//
//	@Summary	Register a new user
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.RegisterRequest	true	"Register request"
//	@Success	200		{object}	models.UserResponse
//	@Router		/auth/register [post]
func (contl *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := contl.AuthService.Register(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.UserResponse{User: user})
}

// Login godoc
//
//	@Summary	Login a user
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.LoginRequest	true	"Login request"
//	@Success	200		{object}	models.LoginResponse
//	@Router		/auth/login [post]
func (contl *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := contl.AuthService.Login(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	session := &models.Session{
		UserID:    user.ID,
		IpAddress: ctx.ClientIP(),
		UserAgent: ctx.Request.UserAgent(),
		ExpiredAt: lib.SESSION_EXPIRED_TIME,
	}
	session, err = contl.SessionService.CreateSession(session)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokens, err := contl.AuthService.GenerateJwtTokens(&models.UserWithClaims{
		User:      *user,
		SessionID: session.ID,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.SetCookie("rt", tokens.RefreshToken, int(lib.REFRESH_TOKEN_TTL_DURATION), "/", "", true, true)
	ctx.JSON(http.StatusOK, &models.LoginResponse{User: user, AccessToken: tokens.AccessToken})
}

func (contl *AuthController) Logout(ctx *gin.Context) {}
