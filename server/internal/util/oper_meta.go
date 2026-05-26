package util

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type OperMeta struct {
	Module       string
	Action       string
	BusinessType string
}

// ResolveOperMeta 根据路由模板与 HTTP 方法生成操作描述。
func ResolveOperMeta(method, fullPath string) OperMeta {
	method = strings.ToUpper(method)
	path := fullPath
	if path == "" {
		path = "/"
	}

	switch {
	case strings.HasPrefix(path, "/api/v1/auth/profile"):
		return OperMeta{Module: "个人中心", Action: actionByMethod(method, "查看资料", "", "修改资料", ""), BusinessType: "profile"}
	case strings.HasPrefix(path, "/api/v1/auth/password"):
		return OperMeta{Module: "个人中心", Action: "修改密码", BusinessType: "profile"}
	case strings.HasPrefix(path, "/api/v1/system/users"):
		return OperMeta{Module: "用户管理", Action: crudAction(method, "用户"), BusinessType: "user"}
	case strings.HasPrefix(path, "/api/v1/system/roles"):
		if strings.Contains(path, "/menus") {
			return OperMeta{Module: "角色管理", Action: actionByMethod(method, "查询角色菜单", "", "分配菜单", ""), BusinessType: "role"}
		}
		return OperMeta{Module: "角色管理", Action: crudAction(method, "角色"), BusinessType: "role"}
	case strings.HasPrefix(path, "/api/v1/system/menus"):
		return OperMeta{Module: "菜单管理", Action: crudAction(method, "菜单"), BusinessType: "menu"}
	case strings.HasPrefix(path, "/api/v1/system/depts"):
		return OperMeta{Module: "部门管理", Action: crudAction(method, "部门"), BusinessType: "dept"}
	case strings.HasPrefix(path, "/api/v1/system/dicts"):
		return OperMeta{Module: "字典管理", Action: crudAction(method, "字典"), BusinessType: "dict"}
	case strings.HasPrefix(path, "/api/v1/system/dict-details"):
		return OperMeta{Module: "字典管理", Action: crudAction(method, "字典项"), BusinessType: "dict_detail"}
	case strings.HasPrefix(path, "/api/v1/system/login-logs"):
		return OperMeta{Module: "登录日志", Action: "查询登录日志", BusinessType: "login_log"}
	case strings.HasPrefix(path, "/api/v1/system/oper-logs"):
		return OperMeta{Module: "操作日志", Action: "查询操作日志", BusinessType: "oper_log"}
	default:
		return OperMeta{Module: "系统", Action: method + " " + path, BusinessType: ""}
	}
}

func crudAction(method, noun string) string {
	return actionByMethod(method, "查询"+noun, "新增"+noun, "修改"+noun, "删除"+noun)
}

func actionByMethod(method string, get, post, put, del string) string {
	switch method {
	case "GET", "HEAD":
		if get != "" {
			return get
		}
		return "查询"
	case "POST":
		if post != "" {
			return post
		}
		return "新增"
	case "PUT", "PATCH":
		if put != "" {
			return put
		}
		return "修改"
	case "DELETE":
		if del != "" {
			return del
		}
		return "删除"
	default:
		return method
	}
}

func ParseBusinessID(c *gin.Context) *int64 {
	if idStr := c.Param("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			return &id
		}
	}
	return nil
}
