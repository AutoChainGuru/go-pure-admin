package router

import (
	"pure-go-admin/server/internal/handler"
	"pure-go-admin/server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.CORS(), middleware.AccessLog())

	authH := handler.AuthHandler{}
	sysH := handler.SystemHandler{}

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authH.Login)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth(), middleware.OperLog())
		{
			auth.GET("/auth/info", authH.Info)
			auth.GET("/auth/profile", authH.GetProfile)
			auth.PUT("/auth/profile", authH.UpdateProfile)
			auth.PUT("/auth/password", authH.ChangePassword)

			sys := auth.Group("/system")
			{
				sys.GET("/users", sysH.ListUsers)
				sys.POST("/users", sysH.CreateUser)
				sys.PUT("/users/:id", sysH.UpdateUser)
				sys.DELETE("/users/:id", sysH.DeleteUser)

				sys.GET("/roles", sysH.ListRoles)
				sys.GET("/roles/all", sysH.AllRoles)
				sys.POST("/roles", sysH.CreateRole)
				sys.PUT("/roles/:id", sysH.UpdateRole)
				sys.DELETE("/roles/:id", sysH.DeleteRole)
				sys.GET("/roles/:id/menus", sysH.GetRoleMenus)
				sys.PUT("/roles/:id/menus", sysH.AssignRoleMenus)

				sys.GET("/menus", sysH.MenuTree)
				sys.POST("/menus", sysH.CreateMenu)
				sys.PUT("/menus/:id", sysH.UpdateMenu)
				sys.DELETE("/menus/:id", sysH.DeleteMenu)

				sys.GET("/depts", sysH.DeptTree)
				sys.POST("/depts", sysH.CreateDept)
				sys.PUT("/depts/:id", sysH.UpdateDept)
				sys.DELETE("/depts/:id", sysH.DeleteDept)

				sys.GET("/dicts", sysH.ListDicts)
				sys.POST("/dicts", sysH.CreateDict)
				sys.PUT("/dicts/:id", sysH.UpdateDict)
				sys.DELETE("/dicts/:id", sysH.DeleteDict)

				sys.GET("/dict-details", sysH.ListDictDetails)
				sys.POST("/dict-details", sysH.CreateDictDetail)
				sys.PUT("/dict-details/:id", sysH.UpdateDictDetail)
				sys.DELETE("/dict-details/:id", sysH.DeleteDictDetail)

				sys.GET("/login-logs", sysH.ListLoginLogs)
				sys.GET("/oper-logs", sysH.ListOperLogs)
				sys.GET("/oper-logs/:id", sysH.GetOperLog)
			}
		}
	}
}
