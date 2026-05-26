package handler

import (
	"strconv"

	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/pkg/response"
	"go-pure-admin/server/internal/service"
	"go-pure-admin/server/internal/util"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	users service.UserService
	roles service.RoleService
	menus service.MenuService
	depts service.DeptService
	dicts service.DictService
	logs  service.LogService
}

// --- users ---

func (h SystemHandler) ListUsers(c *gin.Context) {
	var enabled *bool
	if v := c.Query("enabled"); v != "" {
		b := v == "true" || v == "1"
		enabled = &b
	}
	var deptID *int64
	if v := c.Query("deptId"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		deptID = &id
	}
	list, total, err := h.users.List(service.UserListQuery{
		Username: c.Query("username"),
		Phone:    c.Query("phone"),
		Enabled:  enabled,
		DeptID:   deptID,
	}, util.ParsePage(c))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, list, total, util.ParsePage(c).Page, util.ParsePage(c).PageSize)
}

func (h SystemHandler) CreateUser(c *gin.Context) {
	var req struct {
		model.SysUser
		Password string  `json:"password"`
		RoleIDs  []int64 `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.users.Create(&req.SysUser, req.RoleIDs, req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效ID")
		return
	}
	var req struct {
		Nickname *string `json:"nickname"`
		Phone    *string `json:"phone"`
		Email    *string `json:"email"`
		Avatar   *string `json:"avatar"`
		Enabled  *bool   `json:"enabled"`
		DeptID   *int64  `json:"deptId"`
		Password string  `json:"password"`
		RoleIDs  []int64 `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	fields := map[string]any{}
	if req.Nickname != nil {
		fields["nickname"] = *req.Nickname
	}
	if req.Phone != nil {
		fields["phone"] = *req.Phone
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}
	if req.Avatar != nil {
		fields["avatar"] = *req.Avatar
	}
	if req.Enabled != nil {
		fields["enabled"] = *req.Enabled
	}
	if req.DeptID != nil {
		fields["dept_id"] = *req.DeptID
	}
	if err := h.users.Update(id, fields, req.RoleIDs, req.Password); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.users.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

// --- roles ---

func (h SystemHandler) ListRoles(c *gin.Context) {
	p := util.ParsePage(c)
	list, total, err := h.roles.List(c.Query("keyword"), p)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, list, total, p.Page, p.PageSize)
}

func (h SystemHandler) AllRoles(c *gin.Context) {
	list, err := h.roles.All()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h SystemHandler) CreateRole(c *gin.Context) {
	var req struct {
		model.SysRole
		DeptIDs []int64 `json:"deptIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.roles.Create(&req.SysRole, req.DeptIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		model.SysRole
		DeptIDs []int64 `json:"deptIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.roles.Update(id, &req.SysRole, req.DeptIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.roles.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

func (h SystemHandler) GetRoleMenus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids, err := h.roles.GetMenuIDs(id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"menuIds": ids})
}

func (h SystemHandler) AssignRoleMenus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		MenuIDs            []int64 `json:"menuIds"`
		AutoIncludeButtons *bool   `json:"autoIncludeButtons"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	autoButtons := false
	if req.AutoIncludeButtons != nil {
		autoButtons = *req.AutoIncludeButtons
	}
	if err := h.roles.AssignMenus(id, req.MenuIDs, autoButtons); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "分配成功")
}

// --- menus ---

func (h SystemHandler) MenuTree(c *gin.Context) {
	tree, err := h.menus.Tree(true)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, tree)
}

func (h SystemHandler) CreateMenu(c *gin.Context) {
	var m model.SysMenu
	if err := c.ShouldBindJSON(&m); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.menus.Create(&m); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateMenu(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var m model.SysMenu
	if err := c.ShouldBindJSON(&m); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.menus.Update(id, &m); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteMenu(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.menus.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

// --- depts ---

func (h SystemHandler) DeptTree(c *gin.Context) {
	tree, err := h.depts.Tree()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, tree)
}

func (h SystemHandler) CreateDept(c *gin.Context) {
	var d model.SysDept
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.depts.Create(&d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateDept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var d model.SysDept
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.depts.Update(id, &d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteDept(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.depts.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

// --- dicts ---

func (h SystemHandler) ListDicts(c *gin.Context) {
	p := util.ParsePage(c)
	list, total, err := h.dicts.List(c.Query("keyword"), p)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, list, total, p.Page, p.PageSize)
}

func (h SystemHandler) CreateDict(c *gin.Context) {
	var d model.SysDict
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.dicts.Create(&d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateDict(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var d model.SysDict
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.dicts.Update(id, &d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteDict(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.dicts.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

func (h SystemHandler) ListDictDetails(c *gin.Context) {
	dictID, _ := strconv.ParseInt(c.Query("dictId"), 10, 64)
	list, err := h.dicts.ListDetails(dictID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h SystemHandler) CreateDictDetail(c *gin.Context) {
	var d model.SysDictDetail
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.dicts.CreateDetail(&d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "创建成功")
}

func (h SystemHandler) UpdateDictDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var d model.SysDictDetail
	if err := c.ShouldBindJSON(&d); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.dicts.UpdateDetail(id, &d); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "更新成功")
}

func (h SystemHandler) DeleteDictDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.dicts.DeleteDetail(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMessage(c, "删除成功")
}

// --- logs ---

func (h SystemHandler) ListLoginLogs(c *gin.Context) {
	start, end := util.ParseTimeRange(c)
	p := util.ParsePage(c)
	list, total, err := h.logs.ListLoginLogs(service.LoginLogQuery{
		Username:  c.Query("username"),
		Status:    parseStatusQuery(c.Query("status")),
		IP:        c.Query("ip"),
		StartTime: start,
		EndTime:   end,
	}, p)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, list, total, p.Page, p.PageSize)
}

func (h SystemHandler) ListOperLogs(c *gin.Context) {
	start, end := util.ParseTimeRange(c)
	p := util.ParsePage(c)
	list, total, err := h.logs.ListOperLogs(service.OperLogQuery{
		Username:  c.Query("username"),
		Module:    c.Query("module"),
		Status:    parseStatusQuery(c.Query("status")),
		StartTime: start,
		EndTime:   end,
	}, p)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Page(c, list, total, p.Page, p.PageSize)
}

func (h SystemHandler) GetOperLog(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效ID")
		return
	}
	log, err := h.logs.GetOperLog(id)
	if err != nil {
		response.BadRequest(c, "记录不存在")
		return
	}
	response.OK(c, log)
}

func parseStatusQuery(v string) *int16 {
	if v == "" {
		return nil
	}
	n, err := strconv.ParseInt(v, 10, 16)
	if err != nil {
		return nil
	}
	s := int16(n)
	return &s
}
