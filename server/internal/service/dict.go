package service

import (
	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/util"
)

type DictService struct{}

func (DictService) List(keyword string, p util.Page) ([]model.SysDict, int64, error) {
	db := global.DB.Model(&model.SysDict{})
	if keyword != "" {
		db = db.Where("name ILIKE ? OR type ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SysDict
	err := db.Scopes(util.Paginate(p)).Order("id desc").Find(&list).Error
	return list, total, err
}

func (DictService) Create(d *model.SysDict) error {
	return global.DB.Create(d).Error
}

func (DictService) Update(id int64, d *model.SysDict) error {
	return global.DB.Model(&model.SysDict{}).Where("id = ?", id).Updates(d).Error
}

func (DictService) Delete(id int64) error {
	_ = global.DB.Where("dict_id = ?", id).Delete(&model.SysDictDetail{})
	return global.DB.Delete(&model.SysDict{}, id).Error
}

func (DictService) ListDetails(dictID int64) ([]model.SysDictDetail, error) {
	var list []model.SysDictDetail
	err := global.DB.Where("dict_id = ?", dictID).Order("sort asc, id asc").Find(&list).Error
	return list, err
}

func (DictService) CreateDetail(d *model.SysDictDetail) error {
	return global.DB.Create(d).Error
}

func (DictService) UpdateDetail(id int64, d *model.SysDictDetail) error {
	return global.DB.Model(&model.SysDictDetail{}).Where("id = ?", id).Updates(d).Error
}

func (DictService) DeleteDetail(id int64) error {
	return global.DB.Delete(&model.SysDictDetail{}, id).Error
}
