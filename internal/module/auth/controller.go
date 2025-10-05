package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/middleware"
	"dev-go-apis/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type IAuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.User, error)
	GenerateJwtTokens(userWithClaims *models.UserWithClaims) (*models.JwtTokens, error)
	GenerateOAuthState(prefix string) (string, error)
	CheckOAuthState(state string) (bool, error)
	LoginWithGoogle(user *models.GoogleAccountInfo) (*models.User, error)
	CreateEmailVerification(verification *models.Verification) (*models.CreateVerification, error)
}

type ISessionService interface {
	CreateSession(session *models.Session) (*models.Session, error)
	GetSessionById(id uuid.UUID) (*models.Session, error)
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
	authGroup.GET("/refresh-token", middleware.RefreshTokenHandler(), contl.RefreshToken)
	authGroup.GET("/google", contl.LoginWithGoogle)
	authGroup.GET("/google/callback", contl.GoogleCallback)
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
		response := lib.InvalidBodyRequestError

		if validationErrors := lib.ParseValidationErrors(err); validationErrors != nil {
			response.SetErrors(validationErrors)
			lib.SendErrorResponse(ctx, response)
			return
		}

		lib.SendErrorResponse(ctx, response)
		return
	}

	user, err := contl.AuthService.Login(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	if !user.EmailVerified {
		result, err := contl.AuthService.CreateEmailVerification(&models.Verification{
			UserID:     user.ID,
			Identifier: models.EmailVerifiedIdentifier,
		})
		if err != nil {
			ctx.Error(lib.InternalServerError.WithStack(err.Error()))
			return
		}
		lib.SendEmailVerification(user.Email, user.Name, result.Token)
		lib.SendSucceedResponse(ctx, &models.LoginResponse{
			User:       *user,
			IsVerified: user.EmailVerified,
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
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	tokens, err := contl.AuthService.GenerateJwtTokens(&models.UserWithClaims{
		UserID:    user.ID,
		SessionID: session.ID,
	})
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	ctx.SetCookie("rt", tokens.RefreshToken, int(lib.REFRESH_TOKEN_TTL_DURATION), "/", "", true, true)
	lib.SendSucceedResponse(ctx, &models.LoginResponse{
		User:       *user,
		Token:      tokens.AccessToken,
		IsVerified: user.EmailVerified,
	})
}

// LoginWithGoogle godoc
//
//	@Summary	Log in with Google
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Router		/auth/google [get]
func (contl *AuthController) LoginWithGoogle(ctx *gin.Context) {
	state, err := contl.AuthService.GenerateOAuthState(lib.GoogleStatePrefix)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}
	url := lib.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
//
//	@Summary	Google callback
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.APIResponse
//	@Router		/auth/google/callback [get]
func (contl *AuthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	if _, err := contl.AuthService.CheckOAuthState(state); err != nil {
		lib.SendErrorResponse(ctx, lib.InvalidOAuthStateError.WithStack(err.Error()))
		return
	}

	code := ctx.Query("code")
	token, err := lib.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	client := lib.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get(lib.GoogleOAuthUserInfoURl)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}
	defer resp.Body.Close()

	var googleAccountInfo models.GoogleAccountInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleAccountInfo); err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	user, err := contl.AuthService.LoginWithGoogle(&googleAccountInfo)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, user)
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
		response := lib.InvalidBodyRequestError

		if validationErrors := lib.ParseValidationErrors(err); validationErrors != nil {
			response.SetErrors(validationErrors)
			lib.SendErrorResponse(ctx, response)
			return
		}

		lib.SendErrorResponse(ctx, response)
		return
	}

	user, err := contl.AuthService.Register(&req)
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}

	result, err := contl.AuthService.CreateEmailVerification(&models.Verification{
		UserID:     user.ID,
		Identifier: models.EmailVerifiedIdentifier,
	})
	if err != nil {
		ctx.Error(lib.InternalServerError.WithStack(err.Error()))
		return
	}
	lib.SendWelcomeEmail(user.Email, user.Name, result.Token)
	lib.SendSucceedResponse(ctx, &models.LoginResponse{
		User:       *user,
		IsVerified: user.EmailVerified,
	})
}

// RefreshToken godoc
//
//	@Summary	Refresh token
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	models.APIResponse{data=models.RefreshTokenResponse}
//	@Router		/auth/refresh-token [get]
func (contl *AuthController) RefreshToken(ctx *gin.Context) {
	userWithClaims := ctx.MustGet("user").(*models.UserWithClaims)

	session, err := contl.SessionService.GetSessionById(userWithClaims.SessionID)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	if session == nil {
		lib.SendErrorResponse(ctx, lib.UnauthorizedError)
	}

	tokens, err := contl.AuthService.GenerateJwtTokens(userWithClaims)
	if err != nil {
		lib.SendErrorResponse(ctx, lib.InternalServerError.WithStack(err.Error()))
		return
	}

	lib.SendSucceedResponse(ctx, &models.RefreshTokenResponse{
		Token: tokens.AccessToken,
	})
}
