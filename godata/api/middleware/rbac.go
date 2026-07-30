package middleware

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
)

// RBAC Casbin 权限中间件
func RBAC(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		obj := c.Request.URL.Path
		act := c.Request.Method
		ok, err := enforcer.Enforce(fmt.Sprint(userID), obj, act)
		if err != nil || !ok {
			response.Error(c, errcode.Forbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
