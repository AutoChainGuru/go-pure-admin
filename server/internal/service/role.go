package service

import (
	"errors"

	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/util"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RoleService struct{}

func (RoleService) List(keyword string, p util.Page) ([]model.SysRole, int64, error) {
	db := global.DB.Model(&model.SysRole{})
	if keyword != "" {
		db = db.Where("name ILIKE ? OR code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SysRole
	if err := db.Scopes(util.Paginate(p)).Order("sort asc, id asc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (RoleService) All() ([]model.SysRole, error) {
	var list []model.SysRole
	err := global.DB.Where("enabled = ?", true).Order("sort asc").Find(&list).Error
	return list, err
}

func (RoleService) Create(r *model.SysRole, deptIDs []int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(r).Error; err != nil {
			return err
		}
		return replaceRoleDepts(tx, r.ID, r.DataScope, deptIDs)
	})
}

func (RoleService) Update(id int64, r *model.SysRole, deptIDs []int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SysRole{}).Where("id = ?", id).Updates(r).Error; err != nil {
			return err
		}
		if deptIDs != nil {
			return replaceRoleDepts(tx, id, r.DataScope, deptIDs)
		}
		return nil
	})
}

func (RoleService) Delete(id int64) error {
	var role model.SysRole
	if err := global.DB.First(&role, id).Error; err != nil {
		return err
	}
	if role.Code == "admin" {
		return errors.New("内置管理员角色不可删除")
	}
	var userRef int64
	if err := global.DB.Model(&model.SysUserRole{}).Where("role_id = ?", id).Count(&userRef).Error; err != nil {
		return err
	}
	if userRef > 0 {
		return errors.New("角色已分配给用户，无法删除")
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.SysRoleDept{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SysRole{}, id).Error
	})
}

func (RoleService) GetMenuIDs(roleID int64) ([]int64, error) {
	var ids []int64
	err := global.DB.Model(&model.SysRoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &ids).Error
	return ids, err
}

func (RoleService) AssignMenus(roleID int64, menuIDs []int64, autoIncludeButtons bool) error {
	expanded, err := util.ExpandRoleMenuIDs(global.DB, menuIDs, autoIncludeButtons)
	if err != nil {
		return err
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		for _, mid := range lo.Uniq(expanded) {
			if mid == 0 {
				continue
			}
			if err := tx.Create(&model.SysRoleMenu{RoleID: roleID, MenuID: mid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func replaceRoleDepts(tx *gorm.DB, roleID int64, dataScope int16, deptIDs []int64) error {
	_ = tx.Where("role_id = ?", roleID).Delete(&model.SysRoleDept{})
	if dataScope != 2 {
		return nil
	}
	for _, did := range lo.Uniq(deptIDs) {
		if did == 0 {
			continue
		}
		if err := tx.Create(&model.SysRoleDept{RoleID: roleID, DeptID: did}).Error; err != nil {
			return err
		}
	}
	return nil
}
