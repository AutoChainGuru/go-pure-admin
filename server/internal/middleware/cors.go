package middleware

import (
	"pure-go-admin/server/internal/global"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	origins := global.Cfg.CORS.AllowOrigins
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := ""
		for _, o := range origins {
			if o == origin {
				allowed = origin
				break
			}
		}
		if allowed != "" {
			c.Header("Access-Control-Allow-Origin", allowed)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
