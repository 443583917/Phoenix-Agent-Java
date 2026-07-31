package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/jwt"
	"github.com/phoenix-agent-go/infra/response"
)

// Auth JWT 认证中间件
func Auth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""
		// 优先从 Authorization: Bearer <token> 读取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		}
		// 兼容前端使用的 phoenix-token 请求头
		if tokenStr == "" {
			tokenStr = c.GetHeader("phoenix-token")
		}
		if tokenStr == "" {
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		claims, err := jwtManager.ParseToken(tokenStr)
		if err != nil {
			response.ErrorWithStatus(c, http.StatusUnauthorized, errcode.Unauthorized)
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
