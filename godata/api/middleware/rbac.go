package middleware

import "github.com/gin-gonic/gin"

// RBAC Casbin 权限中间件 — Phase 1 为骨架，Phase 2 实现
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO Phase 2: Casbin enforce
		c.Next()
	}
}
