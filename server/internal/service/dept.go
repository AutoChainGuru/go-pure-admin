package service

import (
	"errors"

	"pure-go-admin/server/internal/global"
	"pure-go-admin/server/internal/model"
	"pure-go-admin/server/internal/util"
)

type DeptService struct{}

func (DeptService) Tree() ([]model.SysDept, error) {
	var list []model.SysDept
	if err := global.DB.Order("sort asc, id asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return util.BuildDeptTree(list, nil), nil
}

func (DeptService) Create(d *model.SysDept) error {
	var parent *model.SysDept
	if d.ParentID != nil {
		var p model.SysDept
		if err := global.DB.First(&p, *d.ParentID).Error; err != nil {
			return err
		}
		parent = &p
	}
	d.Ancestors = util.BuildAncestors(parent)
	return global.DB.Create(d).Error
}

func (DeptService) Update(id int64, d *model.SysDept) error {
	var old model.SysDept
	if err := global.DB.First(&old, id).Error; err != nil {
		return err
	}
	var parent *model.SysDept
	if d.ParentID != nil {
		if *d.ParentID == id {
			return errors.New("父部门不能是自己")
		}
		var p model.SysDept
		if err := global.DB.First(&p, *d.ParentID).Error; err != nil {
			return err
		}
		parent = &p
	}
	d.Ancestors = util.BuildAncestors(parent)
	return global.DB.Model(&model.SysDept{}).Where("id = ?", id).Updates(map[string]any{
		"parent_id": d.ParentID,
		"name":      d.Name,
		"sort":      d.Sort,
		"enabled":   d.Enabled,
		"ancestors": d.Ancestors,
	}).Error
}

func (DeptService) Delete(id int64) error {
	var count int64
	if err := global.DB.Model(&model.SysDept{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("存在子部门，无法删除")
	}
	if err := global.DB.Model(&model.SysUser{}).Where("dept_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("部门下存在用户，无法删除")
	}
	if err := global.DB.Model(&model.SysRoleDept{}).Where("dept_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("部门已被角色数据权限引用，无法删除")
	}
	return global.DB.Delete(&model.SysDept{}, id).Error
}
