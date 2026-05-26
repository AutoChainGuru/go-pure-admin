package global

import (
	"go-pure-admin/server/internal/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Cfg *config.Config
	DB  *gorm.DB
	Log *zap.Logger
)
