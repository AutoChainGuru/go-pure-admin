package middleware

import (
	"strings"

	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/pkg/jwtutil"
	"go-pure-admin/server/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID     = "userId"
	ContextSuperAdmin = "superAdmin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwtutil.Parse(token, global.Cfg.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "登录已失效")
			c.Abort()
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextSuperAdmin, claims.SuperAdmin)
		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	v, _ := c.Get(ContextUserID)
	id, _ := v.(int64)
	return id
}

func IsSuperAdmin(c *gin.Context) bool {
	v, _ := c.Get(ContextSuperAdmin)
	b, _ := v.(bool)
	return b
}
