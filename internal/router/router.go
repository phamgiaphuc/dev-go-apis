package router

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/middleware"
	"dev-go-apis/internal/module/auth"
	"dev-go-apis/internal/module/cache"
	"dev-go-apis/internal/module/permission"
	"dev-go-apis/internal/module/role"
	"dev-go-apis/internal/module/session"
	"dev-go-apis/internal/module/user"
	"dev-go-apis/internal/views"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type IController interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

type Router struct {
	DBClient    *sqlx.DB
	CacheClient *redis.Client
	Router      *gin.Engine
}

func NewRouter(dbClient *sqlx.DB, cacheClient *redis.Client) *Router {
	gin.SetMode(lib.GIN_MODE)
	router := gin.Default()

	router.Use(middleware.ErrorsHandler())

	return &Router{
		DBClient:    dbClient,
		CacheClient: cacheClient,
		Router:      router,
	}
}

func (r *Router) InitRoutes() http.Handler {
	r.Router.Static("/static", "internal/static/")
	r.Router.StaticFile("/styles.css", "internal/views/styles.css")
	r.Router.StaticFile("/favicon.ico", "internal/static/favicon.ico")

	r.Router.GET("/", func(ctx *gin.Context) {
		views.Index().Render(ctx.Request.Context(), ctx.Writer)
	})

	r.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	/**
	 * Middlewares
	 */
	r.Router.Use(middleware.ApiKeyHandler())

	apiGroup := r.Router.Group("/api")

	/**
	 * Repositories
	 */
	cacheRepo := cache.NewCacheRepository(r.CacheClient)
	userRepo := user.NewUserRepository(r.DBClient)
	sessionRepo := session.NewSessionRepository(r.DBClient)
	permissionRepo := permission.NewPermissionRepository(r.DBClient)
	roleRepo := role.NewRoleRepository(r.DBClient)

	/**
	 * Services
	 */
	cacheService := cache.NewCacheService(cacheRepo)
	authService := auth.NewAuthService(userRepo)
	userService := user.NewUserService(userRepo)
	sessionService := session.NewSessionService(sessionRepo)
	permissionService := permission.NewPermissionService(permissionRepo)
	roleService := role.NewRoleService(roleRepo, cacheRepo)

	/**
	 * Controllers
	 */
	routes := []IController{
		auth.NewAuthController(authService, sessionService),
		user.NewUserController(userService, cacheService),
		permission.NewPermissionController(permissionService),
		role.NewRoleController(roleService),
	}
	for _, route := range routes {
		route.RegisterRoutes(apiGroup)
	}

	r.Router.NoRoute(func(ctx *gin.Context) {
		lib.SendErrorResponse(ctx, lib.ResourceNotFoundError)
	})

	return r.Router.Handler()
}
