package main

import (
	"fmt"
	"os"
	"path/filepath"

	"pure-go-admin/server/internal/config"
	"pure-go-admin/server/internal/database"
	"pure-go-admin/server/internal/global"
	"pure-go-admin/server/internal/model"
	"pure-go-admin/server/internal/pkg/logger"
	"pure-go-admin/server/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = resolveConfigPath("config.yaml")
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}
	global.Cfg = cfg

	if err := logger.Init(cfg.Log); err != nil {
		panic(err)
	}
	global.Log = logger.L
	defer logger.Sync()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		global.Log.Fatal("database connect failed", zap.Error(err))
	}
	global.DB = db
	if err := database.AutoMigrate(db); err != nil {
		global.Log.Fatal("auto migrate failed", zap.Error(err))
	}
	if err := model.Seed(db); err != nil {
		global.Log.Fatal("seed failed", zap.Error(err))
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	router.Setup(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	global.Log.Info("server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		global.Log.Fatal("server stopped", zap.Error(err))
	}
}

func resolveConfigPath(name string) string {
	candidates := []string{
		name,
		filepath.Join("server", name),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return name
}
