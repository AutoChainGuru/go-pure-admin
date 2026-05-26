package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/service"
	"go-pure-admin/server/internal/util"

	"github.com/gin-gonic/gin"
)

var operLogSvc = service.LogService{}

// OperLog 记录已登录用户的写操作（POST/PUT/PATCH/DELETE）。
func OperLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldRecordOper(c.Request.Method, c.FullPath()) {
			c.Next()
			return
		}

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		start := time.Now()
		c.Next()

		userID := GetUserID(c)
		username := operLogSvc.FindUsername(userID)
		meta := util.ResolveOperMeta(c.Request.Method, c.FullPath())
		status := model.LogStatusSuccess
		errMsg := ""
		if c.Writer.Status() >= 400 {
			status = model.LogStatusFail
			if msg, ok := c.Get("operLogError"); ok {
				errMsg, _ = msg.(string)
			}
			if errMsg == "" {
				errMsg = c.Errors.String()
			}
			if errMsg == "" {
				errMsg = "请求失败"
			}
		}

		entry := &model.SysOperLog{
			UserID:       ptrInt64(userID),
			Username:     username,
			Module:       meta.Module,
			Action:       meta.Action,
			Method:       c.Request.Method,
			Path:         c.FullPath(),
			BusinessType: meta.BusinessType,
			BusinessID:   util.ParseBusinessID(c),
			Status:       status,
			ErrorMsg:     trimErr(errMsg),
			RequestBody:  util.SanitizeLogBody(string(bodyBytes)),
			DurationMs:   int(time.Since(start).Milliseconds()),
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
		}
		go operLogSvc.RecordOper(entry)
	}
}

func shouldRecordOper(method, path string) bool {
	method = strings.ToUpper(method)
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
	default:
		return false
	}
	if path == "/api/v1/auth/login" {
		return false
	}
	return strings.HasPrefix(path, "/api/v1/")
}

func ptrInt64(v int64) *int64 {
	if v <= 0 {
		return nil
	}
	return &v
}

func trimErr(s string) string {
	if len(s) > 500 {
		return s[:500]
	}
	return s
}
