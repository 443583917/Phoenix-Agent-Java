package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/api/middleware"
	"github.com/phoenix-agent-go/infra/config"
	"github.com/phoenix-agent-go/infra/response"
)

func SetupRouter(cfg *config.AppConfig) *gin.Engine {
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
	{
		auth := r.Group("/api/privilege/auth")
		auth.Use(middleware.RateLimit())
		// Phase 2: auth.POST("/login", handler.Login)
		// Phase 2: auth.POST("/logout", handler.Logout)
		// Phase 2: auth.GET("/captcha", handler.Captcha)
		_ = auth
	}

	// API 路由（需 JWT）
	api := r.Group("/api")
	api.Use(middleware.Auth())
	api.Use(middleware.RBAC())
	{
		// Phase 4: agentGroup := api.Group("/agent")
		// Phase 5: datasourceGroup := api.Group("/datasource")
		// Phase 5: chatGroup := api.Group("")
		_ = api
	}

	// 平台管理路由
	platform := r.Group("/platform")
	platform.Use(middleware.Auth())
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
