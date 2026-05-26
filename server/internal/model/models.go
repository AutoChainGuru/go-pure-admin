package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SysDept struct {
	Base
	ParentID  *int64    `json:"parentId"`
	Ancestors string    `json:"ancestors" gorm:"size:512;not null;default:''"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	Sort      int       `json:"sort" gorm:"not null;default:0"`
	Enabled   bool      `json:"enabled" gorm:"not null;default:true"`
	Children  []SysDept `json:"children,omitempty" gorm:"-"`
}

func (SysDept) TableName() string { return "sys_dept" }

type SysUser struct {
	Base
	DeptID      *int64     `json:"deptId"`
	Username    string     `json:"username" gorm:"size:64;not null"`
	Password    string     `json:"-" gorm:"size:255;not null"`
	Nickname    string     `json:"nickname" gorm:"size:64"`
	Phone       string     `json:"phone" gorm:"size:32"`
	Email       string     `json:"email" gorm:"size:128"`
	Avatar      string     `json:"avatar" gorm:"size:512"`
	Enabled     bool       `json:"enabled" gorm:"not null;default:true"`
	SuperAdmin  bool       `json:"superAdmin" gorm:"not null;default:false"`
	LastLoginAt *time.Time `json:"lastLoginAt"`
	RoleIDs     []int64    `json:"roleIds,omitempty" gorm:"-"`
	Dept        *SysDept   `json:"dept,omitempty" gorm:"foreignKey:DeptID"`
	Roles       []SysRole  `json:"roles,omitempty" gorm:"many2many:sys_user_roles;joinForeignKey:user_id;joinReferences:role_id"`
}

func (SysUser) TableName() string { return "sys_users" }

type SysRole struct {
	Base
	Name        string    `json:"name" gorm:"size:128;not null"`
	Code        string    `json:"code" gorm:"size:64;not null"`
	Sort        int       `json:"sort" gorm:"not null;default:0"`
	Enabled     bool      `json:"enabled" gorm:"not null;default:true"`
	DataScope   int16     `json:"dataScope" gorm:"not null;default:5"`
	Description string    `json:"description" gorm:"size:512"`
	MenuIDs     []int64   `json:"menuIds,omitempty" gorm:"-"`
	DeptIDs     []int64   `json:"deptIds,omitempty" gorm:"-"`
	Menus       []SysMenu `json:"menus,omitempty" gorm:"many2many:sys_role_menus;joinForeignKey:role_id;joinReferences:menu_id"`
}

func (SysRole) TableName() string { return "sys_roles" }

type SysUserRole struct {
	UserID int64 `gorm:"primaryKey;column:user_id"`
	RoleID int64 `gorm:"primaryKey;column:role_id"`
}

func (SysUserRole) TableName() string { return "sys_user_roles" }

type SysRoleDept struct {
	RoleID int64 `gorm:"primaryKey;column:role_id"`
	DeptID int64 `gorm:"primaryKey;column:dept_id"`
}

func (SysRoleDept) TableName() string { return "sys_role_depts" }

type SysMenu struct {
	Base
	ParentID  *int64    `json:"parentId"`
	MenuType  int16     `json:"menuType" gorm:"not null;default:1"`
	PermCode  string    `json:"permCode" gorm:"size:128"`
	Path      string    `json:"path" gorm:"size:255"`
	Name      string    `json:"name" gorm:"size:128"`
	Component string    `json:"component" gorm:"size:255"`
	Redirect  string    `json:"redirect" gorm:"size:255"`
	Title     string    `json:"title" gorm:"size:128;not null"`
	Icon      string    `json:"icon" gorm:"size:128"`
	Sort      int       `json:"sort" gorm:"not null;default:0"`
	Hidden    bool      `json:"hidden" gorm:"not null;default:false"`
	KeepAlive bool      `json:"keepAlive" gorm:"not null;default:false"`
	Enabled   bool      `json:"enabled" gorm:"not null;default:true"`
	Children  []SysMenu `json:"children,omitempty" gorm:"-"`
}

func (SysMenu) TableName() string { return "sys_menus" }

type SysRoleMenu struct {
	RoleID int64 `gorm:"primaryKey;column:role_id"`
	MenuID int64 `gorm:"primaryKey;column:menu_id"`
}

func (SysRoleMenu) TableName() string { return "sys_role_menus" }

type SysDict struct {
	Base
	Name        string `json:"name" gorm:"size:128;not null"`
	Type        string `json:"type" gorm:"size:64;not null"`
	Description string `json:"description" gorm:"size:512"`
	Enabled     bool   `json:"enabled" gorm:"not null;default:true"`
}

func (SysDict) TableName() string { return "sys_dict" }

type SysDictDetail struct {
	Base
	DictID    int64  `json:"dictId" gorm:"not null"`
	Label     string `json:"label" gorm:"size:128;not null"`
	Value     string `json:"value" gorm:"size:128;not null"`
	Sort      int    `json:"sort" gorm:"not null;default:0"`
	Enabled   bool   `json:"enabled" gorm:"not null;default:true"`
	CssClass  string `json:"cssClass" gorm:"size:128"`
	ListClass string `json:"listClass" gorm:"size:128"`
	IsDefault bool   `json:"isDefault" gorm:"not null;default:false"`
	Remark    string `json:"remark" gorm:"size:512"`
}

func (SysDictDetail) TableName() string { return "sys_dict_detail" }

func AllModels() []any {
	return []any{
		&SysDept{}, &SysUser{}, &SysRole{}, &SysUserRole{}, &SysRoleDept{},
		&SysMenu{}, &SysRoleMenu{},
		&SysDict{}, &SysDictDetail{},
		&SysLoginLog{}, &SysOperLog{},
	}
}
