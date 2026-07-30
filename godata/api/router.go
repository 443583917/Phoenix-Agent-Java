package api

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/api/handler/privilege"
	"github.com/phoenix-agent-go/api/middleware"
	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
	"github.com/phoenix-agent-go/internal/service"
)

func SetupRouter(cfg *config.AppConfig, jwtManager *jwt.JWTManager, enforcer *casbin.Enforcer, privilegeSvc *service.PrivilegeService) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Tracing())
	r.Use(middleware.CORS(&cfg.Cors))

	// 健康检查
	r.GET("/echo", func(c *gin.Context) {
		response.Success(c, "ok")
	})

	// 静态文件（头像等）
	r.Static("/api/upload", "./storage/upload")

	// 认证路由（无需 JWT）
	authHandler := privilege.NewAuthHandler(privilegeSvc, jwtManager)
	auth := r.Group("/api/privilege/auth")
	auth.Use(middleware.RateLimit())
	auth.GET("/captcha", authHandler.Captcha)
	auth.POST("/login", authHandler.Login)
	auth.POST("/logout", authHandler.Logout)
	auth.GET("/menus", middleware.Auth(jwtManager), authHandler.Menus)
	auth.GET("/getLoginUserInfo", middleware.Auth(jwtManager), authHandler.GetLoginUserInfo)

	// API 路由（需 JWT）
	api := r.Group("/api")
	api.Use(middleware.Auth(jwtManager))
	api.Use(middleware.RBAC(enforcer))
	{
		// 用户管理
		userHandler := privilege.NewUserHandler(privilegeSvc)
		userGroup := api.Group("/privilege/user")
		{
			userGroup.GET("/page", userHandler.Page)
			userGroup.GET("/:id", userHandler.GetByID)
			userGroup.GET("/code/:code", userHandler.GetByCode)
			userGroup.POST("", userHandler.Create)
			userGroup.PUT("", userHandler.Update)
			userGroup.DELETE("/:id", userHandler.Delete)
			userGroup.PUT("/password", userHandler.UpdatePassword)
			userGroup.PUT("/reset-password/:id", userHandler.ResetPassword)
		}

		// 角色管理
		roleHandler := privilege.NewRoleHandler(privilegeSvc)
		roleGroup := api.Group("/privilege/role")
		{
			roleGroup.POST("/page", roleHandler.Page)
			roleGroup.GET("/:id", roleHandler.GetByID)
			roleGroup.GET("/company/:companyId", roleHandler.GetByCompany)
			roleGroup.POST("", roleHandler.Create)
			roleGroup.PUT("", roleHandler.Update)
			roleGroup.DELETE("/:id", roleHandler.Delete)
			roleGroup.GET("/:roleId/acls", roleHandler.GetAcls)
		}

		// 用户-角色关联管理
		userRoleHandler := privilege.NewUserRoleHandler(privilegeSvc)
		userRoleGroup := api.Group("/privilege/user-role")
		{
			userRoleGroup.GET("/user/:userId", userRoleHandler.GetByUser)
			userRoleGroup.GET("/role/:roleId", userRoleHandler.GetByRole)
			userRoleGroup.POST("", userRoleHandler.Create)
			userRoleGroup.DELETE("/:id", userRoleHandler.Delete)
			userRoleGroup.POST("/batch-save", userRoleHandler.BatchSave)
			userRoleGroup.DELETE("/batch-remove", userRoleHandler.BatchRemove)
		}

		// Phase 4: agentGroup := api.Group("/agent")
		// Phase 5: datasourceGroup := api.Group("/datasource")
		// Phase 5: chatGroup := api.Group("")
	}

	// 平台管理路由
	platform := r.Group("/platform")
	platform.Use(middleware.Auth(jwtManager))
	{
		// Phase 3: platform handler registration
		_ = platform
	}

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.Response{
			Code:    404,
			Message: "not found",
		})
	})

	return r
}
