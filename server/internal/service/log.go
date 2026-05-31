package service

import (
	"time"

	"pure-go-admin/server/internal/global"
	"pure-go-admin/server/internal/model"
	"pure-go-admin/server/internal/util"
)

type LogService struct{}

type LoginLogQuery struct {
	Username  string
	Status    *int16
	IP        string
	StartTime *time.Time
	EndTime   *time.Time
}

type OperLogQuery struct {
	Username  string
	Module    string
	Status    *int16
	StartTime *time.Time
	EndTime   *time.Time
}

func (LogService) RecordLogin(entry *model.SysLoginLog) {
	if entry == nil {
		return
	}
	_ = global.DB.Create(entry).Error
}

func (LogService) RecordOper(entry *model.SysOperLog) {
	if entry == nil {
		return
	}
	_ = global.DB.Create(entry).Error
}

func (LogService) ListLoginLogs(q LoginLogQuery, p util.Page) ([]model.SysLoginLog, int64, error) {
	db := global.DB.Model(&model.SysLoginLog{})
	if q.Username != "" {
		db = db.Where("username ILIKE ?", "%"+q.Username+"%")
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.IP != "" {
		db = db.Where("ip ILIKE ?", "%"+q.IP+"%")
	}
	if q.StartTime != nil {
		db = db.Where("created_at >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		db = db.Where("created_at <= ?", *q.EndTime)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SysLoginLog
	err := db.Scopes(util.Paginate(p)).Order("id desc").Find(&list).Error
	return list, total, err
}

func (LogService) ListOperLogs(q OperLogQuery, p util.Page) ([]model.SysOperLog, int64, error) {
	db := global.DB.Model(&model.SysOperLog{})
	if q.Username != "" {
		db = db.Where("username ILIKE ?", "%"+q.Username+"%")
	}
	if q.Module != "" {
		db = db.Where("module ILIKE ?", "%"+q.Module+"%")
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.StartTime != nil {
		db = db.Where("created_at >= ?", *q.StartTime)
	}
	if q.EndTime != nil {
		db = db.Where("created_at <= ?", *q.EndTime)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.SysOperLog
	err := db.Scopes(util.Paginate(p)).Order("id desc").Find(&list).Error
	return list, total, err
}

func (LogService) GetOperLog(id int64) (*model.SysOperLog, error) {
	var log model.SysOperLog
	if err := global.DB.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (LogService) FindUsername(userID int64) string {
	if userID <= 0 {
		return ""
	}
	var u model.SysUser
	if err := global.DB.Select("username").First(&u, userID).Error; err != nil {
		return ""
	}
	return u.Username
}
