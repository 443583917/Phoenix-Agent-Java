package middleware

import "github.com/gin-gonic/gin"

// Auth JWT + OAuth2 认证中间件 — Phase 1 为骨架，Phase 2 实现
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO Phase 2: JWT 验证 + Casbin 加载角色
		c.Next()
	}
}
