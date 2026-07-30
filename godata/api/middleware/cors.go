package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/internal/config"
)

func CORS(cfg *config.CorsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Next()
			return
		}

		allowOrigin := ""
		for _, o := range cfg.AllowOrigins {
			if o == "*" {
				if cfg.AllowCredentials {
					allowOrigin = origin // must echo specific origin when credentials enabled
				} else {
					allowOrigin = "*"
				}
				break
			}
			if o == origin {
				allowOrigin = origin
				break
			}
		}

		if allowOrigin == "" {
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(cfg.AllowCredentials))
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,X-Trace-Id")
		c.Header("Access-Control-Expose-Headers", "Content-Length,X-Trace-Id")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
