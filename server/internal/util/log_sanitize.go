package util

import (
	"encoding/json"
	"strings"
)

const maxLogBodyLen = 4000

var sensitiveKeys = []string{"password", "oldPassword", "newPassword", "confirmPassword", "token"}

// SanitizeLogBody 脱敏并截断请求/响应体。
func SanitizeLogBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	var data any
	if err := json.Unmarshal([]byte(body), &data); err == nil {
		maskSensitive(data)
		if b, err := json.Marshal(data); err == nil {
			body = string(b)
		}
	} else {
		lower := strings.ToLower(body)
		for _, k := range sensitiveKeys {
			if strings.Contains(lower, k) {
				return "[已脱敏]"
			}
		}
	}
	if len(body) > maxLogBodyLen {
		return body[:maxLogBodyLen] + "…"
	}
	return body
}

func maskSensitive(v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if isSensitiveKey(k) {
				t[k] = "***"
				continue
			}
			maskSensitive(val)
		}
	case []any:
		for _, item := range t {
			maskSensitive(item)
		}
	}
}

func isSensitiveKey(k string) bool {
	kl := strings.ToLower(k)
	for _, s := range sensitiveKeys {
		if kl == s {
			return true
		}
	}
	return false
}
