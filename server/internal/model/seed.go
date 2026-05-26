package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&SysUser{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	dept := SysDept{Name: "总公司", Ancestors: "0", Sort: 0, Enabled: true}
	if err := db.Create(&dept).Error; err != nil {
		return err
	}

	role := SysRole{
		Name: "超级管理员", Code: "admin", Sort: 0, Enabled: true,
		DataScope: 1, Description: "拥有全部权限",
	}
	if err := db.Create(&role).Error; err != nil {
		return err
	}

	admin := SysUser{
		DeptID: &dept.ID, Username: "admin", Password: string(hash),
		Nickname: "管理员", Enabled: true, SuperAdmin: true,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&SysUserRole{UserID: admin.ID, RoleID: role.ID}).Error; err != nil {
		return err
	}

	menus, err := seedMenus(db)
	if err != nil {
		return err
	}

	for _, m := range menus {
		if err := db.Create(&SysRoleMenu{RoleID: role.ID, MenuID: m.ID}).Error; err != nil {
			return err
		}
	}

	dict := SysDict{Name: "通用状态", Type: "sys_status", Description: "启用禁用", Enabled: true}
	if err := db.Create(&dict).Error; err != nil {
		return err
	}
	for i, item := range []struct{ label, value string }{
		{"启用", "1"}, {"禁用", "0"},
	} {
		if err := db.Create(&SysDictDetail{
			DictID: dict.ID, Label: item.label, Value: item.value, Sort: i, Enabled: true,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedMenus(db *gorm.DB) ([]SysMenu, error) {
	type menuDef struct {
		title, path, name, component, icon string
		menuType                           int16
		sort                               int
		permCode                           string
		children                           []menuDef
		buttons                            []menuDef
	}

	tree := []menuDef{
		{
			title: "工作台", path: "/dashboard", name: "Dashboard", component: "dashboard/index",
			menuType: 1, sort: 0, icon: "Odometer",
		},
		{
			title: "系统管理", path: "/system", name: "System", component: "Layout",
			menuType: 0, sort: 1, icon: "Setting",
			children: []menuDef{
				{title: "用户管理", path: "user", name: "SystemUser", component: "system/user/index", menuType: 1, sort: 1, icon: "User",
					buttons: []menuDef{
						{title: "新增用户", permCode: "system:user:add", menuType: 2},
						{title: "编辑用户", permCode: "system:user:edit", menuType: 2},
						{title: "删除用户", permCode: "system:user:delete", menuType: 2},
					}},
				{title: "角色管理", path: "role", name: "SystemRole", component: "system/role/index", menuType: 1, sort: 2, icon: "UserFilled",
					buttons: []menuDef{
						{title: "新增角色", permCode: "system:role:add", menuType: 2},
						{title: "编辑角色", permCode: "system:role:edit", menuType: 2},
						{title: "删除角色", permCode: "system:role:delete", menuType: 2},
					}},
				{title: "菜单管理", path: "menu", name: "SystemMenu", component: "system/menu/index", menuType: 1, sort: 3, icon: "Menu"},
				{title: "部门管理", path: "dept", name: "SystemDept", component: "system/dept/index", menuType: 1, sort: 4, icon: "OfficeBuilding"},
				{title: "字典管理", path: "dict", name: "SystemDict", component: "system/dict/index", menuType: 1, sort: 5, icon: "Collection"},
				{title: "登录日志", path: "login-log", name: "SystemLoginLog", component: "system/login-log/index", menuType: 1, sort: 6, icon: "Key"},
				{title: "操作日志", path: "oper-log", name: "SystemOperLog", component: "system/oper-log/index", menuType: 1, sort: 7, icon: "Document"},
			},
		},
	}

	var allMenus []SysMenu
	var createMenu func(def menuDef, parentID *int64) error
	createMenu = func(def menuDef, parentID *int64) error {
		m := SysMenu{
			ParentID: parentID, MenuType: def.menuType, PermCode: def.permCode,
			Path: def.path, Name: def.name, Component: def.component, Title: def.title,
			Icon: def.icon, Sort: def.sort, Enabled: true,
		}
		if err := db.Create(&m).Error; err != nil {
			return err
		}
		allMenus = append(allMenus, m)
		for _, b := range def.buttons {
			if err := createMenu(b, &m.ID); err != nil {
				return err
			}
		}
		for _, c := range def.children {
			if err := createMenu(c, &m.ID); err != nil {
				return err
			}
		}
		return nil
	}

	for _, root := range tree {
		if err := createMenu(root, nil); err != nil {
			return nil, err
		}
	}

	return allMenus, nil
}
