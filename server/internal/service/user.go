package service

import (
	"errors"

	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/util"

	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

type UserListQuery struct {
	Username string
	Phone    string
	Enabled  *bool
	DeptID   *int64
}

func (UserService) List(q UserListQuery, p util.Page) ([]model.SysUser, int64, error) {
	db := global.DB.Model(&model.SysUser{}).Preload("Dept").Preload("Roles")
	if q.Username != "" {
		db = db.Where("username ILIKE ?", "%"+q.Username+"%")
	}
	if q.Phone != "" {
		db = db.Where("phone ILIKE ?", "%"+q.Phone+"%")
	}
	if q.Enabled != nil {
		db = db.Where("enabled = ?", *q.Enabled)
	}
	if q.DeptID != nil {
		db = db.Where("dept_id = ?", *q.DeptID)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SysUser
	if err := db.Scopes(util.Paginate(p)).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].Password = ""
	}
	return list, total, nil
}

func (UserService) Create(u *model.SysUser, roleIDs []int64, password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		if err := replaceUserRoles(tx, u.ID, roleIDs); err != nil {
			return err
		}
		return nil
	})
}

func (UserService) Update(id int64, fields map[string]any, roleIDs []int64, password string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			fields["password"] = string(hash)
		}
		if len(fields) > 0 {
			if err := tx.Model(&model.SysUser{}).Where("id = ?", id).Updates(fields).Error; err != nil {
				return err
			}
		}
		if roleIDs != nil {
			if err := replaceUserRoles(tx, id, roleIDs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (UserService) Delete(id int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.SysUser{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func replaceUserRoles(tx *gorm.DB, userID int64, roleIDs []int64) error {
	if err := tx.Where("user_id = ?", userID).Delete(&model.SysUserRole{}).Error; err != nil {
		return err
	}
	roleIDs = lo.Uniq(roleIDs)
	for _, rid := range roleIDs {
		if rid == 0 {
			continue
		}
		if err := tx.Create(&model.SysUserRole{UserID: userID, RoleID: rid}).Error; err != nil {
			return err
		}
	}
	return nil
}
