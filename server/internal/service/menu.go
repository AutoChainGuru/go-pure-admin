package service

import (
	"errors"
	"strings"

	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/util"
)

func validateMenu(m *model.SysMenu) error {
	if strings.TrimSpace(m.Title) == "" {
		return errors.New("标题不能为空")
	}
	switch m.MenuType {
	case 2:
		if strings.TrimSpace(m.PermCode) == "" {
			return errors.New("权限码不能为空")
		}
	case 0, 1:
		if strings.TrimSpace(m.Path) == "" {
			return errors.New("路径不能为空")
		}
		if strings.TrimSpace(m.Name) == "" {
			return errors.New("路由名不能为空")
		}
		if strings.TrimSpace(m.Component) == "" {
			return errors.New("组件不能为空")
		}
	default:
		return errors.New("无效的菜单类型")
	}
	return nil
}

type MenuService struct{}

func (MenuService) Tree(all bool) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	q := global.DB.Order("sort asc, id asc")
	if !all {
		q = q.Where("enabled = ?", true)
	}
	if err := q.Find(&menus).Error; err != nil {
		return nil, err
	}
	return util.BuildMenuTree(menus, nil), nil
}

func (MenuService) Create(m *model.SysMenu) error {
	if err := validateMenu(m); err != nil {
		return err
	}
	return global.DB.Create(m).Error
}

func (MenuService) Update(id int64, m *model.SysMenu) error {
	if err := validateMenu(m); err != nil {
		return err
	}
	updates := map[string]any{
		"menu_type":  m.MenuType,
		"perm_code":  m.PermCode,
		"path":       m.Path,
		"name":       m.Name,
		"component":  m.Component,
		"redirect":   m.Redirect,
		"title":      m.Title,
		"icon":       m.Icon,
		"sort":       m.Sort,
		"hidden":     m.Hidden,
		"keep_alive": m.KeepAlive,
		"enabled":    m.Enabled,
		"parent_id":  m.ParentID,
	}
	return global.DB.Model(&model.SysMenu{}).Where("id = ?", id).Updates(updates).Error
}

func (MenuService) Delete(id int64) error {
	var childCount int64
	if err := global.DB.Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在子菜单，无法删除")
	}
	var roleRef int64
	if err := global.DB.Model(&model.SysRoleMenu{}).Where("menu_id = ?", id).Count(&roleRef).Error; err != nil {
		return err
	}
	if roleRef > 0 {
		return errors.New("菜单已被角色引用，无法删除")
	}
	return global.DB.Delete(&model.SysMenu{}, id).Error
}
