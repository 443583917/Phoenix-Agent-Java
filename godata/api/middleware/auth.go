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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		claims, err := jwtManager.ParseToken(parts[1])
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
