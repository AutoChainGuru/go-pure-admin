package util

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ParseTimeRange 解析查询参数 startTime、endTime（支持 RFC3339、2006-01-02 15:04:05、2006-01-02）。
func ParseTimeRange(c *gin.Context) (start, end *time.Time) {
	start = parseTimeParam(c.Query("startTime"), false)
	end = parseTimeParam(c.Query("endTime"), true)
	return start, end
}

func parseTimeParam(v string, endOfDay bool) *time.Time {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, v, time.Local); err == nil {
			if layout == "2006-01-02" && endOfDay {
				t = t.Add(24*time.Hour - time.Nanosecond)
			}
			return &t
		}
	}
	return nil
}
