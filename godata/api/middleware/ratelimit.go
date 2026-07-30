package middleware

import "github.com/gin-gonic/gin"

// RateLimit 限流中间件 — Phase 1 为骨架，后续 Phase 实现
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
