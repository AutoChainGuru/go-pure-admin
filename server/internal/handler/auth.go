package handler

import (
	"pure-go-admin/server/internal/middleware"
	"pure-go-admin/server/internal/pkg/response"
	"pure-go-admin/server/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc service.AuthService
}

func (h AuthHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	res, err := h.svc.Login(req, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, res)
}

func (h AuthHandler) Info(c *gin.Context) {
	info, err := h.svc.Info(middleware.GetUserID(c), middleware.IsSuperAdmin(c))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, info)
}

func (h AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.svc.GetProfile(middleware.GetUserID(c))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, user)
}

func (h AuthHandler) UpdateProfile(c *gin.Context) {
	var req service.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	user, err := h.svc.UpdateProfile(middleware.GetUserID(c), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, user)
}

func (h AuthHandler) ChangePassword(c *gin.Context) {
	var req service.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.ChangePassword(middleware.GetUserID(c), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "密码已修改")
}
