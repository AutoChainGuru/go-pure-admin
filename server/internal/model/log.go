package model

import "time"

const (
	LogStatusFail    int16 = 0
	LogStatusSuccess int16 = 1
)

type SysLoginLog struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserID    *int64    `json:"userId"`
	Username  string    `json:"username" gorm:"size:64;not null"`
	Nickname  string    `json:"nickname" gorm:"size:64"`
	Status    int16     `json:"status" gorm:"not null;default:0"`
	Message   string    `json:"message" gorm:"size:512;not null;default:''"`
	IP        string    `json:"ip" gorm:"size:64"`
	Location  string    `json:"location" gorm:"size:255"`
	UserAgent string    `json:"userAgent" gorm:"size:512"`
	CreatedAt time.Time `json:"createdAt"`
}

func (SysLoginLog) TableName() string { return "sys_login_log" }

type SysOperLog struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	UserID       *int64    `json:"userId"`
	Username     string    `json:"username" gorm:"size:64"`
	Module       string    `json:"module" gorm:"size:128;not null;default:''"`
	Action       string    `json:"action" gorm:"size:128;not null;default:''"`
	Method       string    `json:"method" gorm:"size:16"`
	Path         string    `json:"path" gorm:"size:255"`
	BusinessType string    `json:"businessType" gorm:"size:64"`
	BusinessID   *int64    `json:"businessId"`
	Status       int16     `json:"status" gorm:"not null;default:1"`
	ErrorMsg     string    `json:"errorMsg" gorm:"size:512;not null;default:''"`
	RequestBody  string    `json:"requestBody" gorm:"type:text"`
	ResponseBody string    `json:"responseBody" gorm:"type:text"`
	DurationMs   int       `json:"durationMs" gorm:"not null;default:0"`
	IP           string    `json:"ip" gorm:"size:64"`
	UserAgent    string    `json:"userAgent" gorm:"size:512"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (SysOperLog) TableName() string { return "sys_oper_log" }
