package auth

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"log"

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
}

// Register godoc
//
//	@Summary	Register a new user
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.RegisterRequest	true	"Register request"
//	@Success	200		{object}	models.APIResponse{data=models.RegisterResponse}
//	@Failure	400		{object}	models.APIResponse
//	@Router		/auth/register [post]
func (contl *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Binding failed: %v\n", err)
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	user, err := contl.AuthService.Register(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, &models.RegisterResponse{User: user})
}

// Login godoc
//
//	@Summary	Login a user
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		models.LoginRequest	true	"Login request"
//	@Success	200		{object}	models.APIResponse{data=models.LoginResponse}
//	@Failure	400		{object}	models.APIResponse
//	@Router		/auth/login [post]
func (contl *AuthController) Login(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(lib.InvalidBodyRequestError.WithStack(err.Error()))
		return
	}

	user, err := contl.AuthService.Login(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
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
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	tokens, err := contl.AuthService.GenerateJwtTokens(&models.UserWithClaims{
		User:      *user,
		SessionID: session.ID,
	})
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	ctx.SetCookie("rt", tokens.RefreshToken, int(lib.REFRESH_TOKEN_TTL_DURATION), "/", "", true, true)
	lib.SendSucceedResponse(ctx, &models.LoginResponse{User: user, AccessToken: tokens.AccessToken})
}
