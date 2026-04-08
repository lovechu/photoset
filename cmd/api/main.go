package main

import (
	"log"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/http/routes"
	"photoset/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 初始化 MySQL
	if err := database.InitMySQL(cfg); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	defer database.CloseMySQL()

	// 自动建表
	if err := database.GetMySQL().AutoMigrate(
		&domain.User{},
		&domain.PhotoSet{},
		&domain.Photo{},
		&domain.Tag{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// 初始化 Redis
	if err := database.InitRedis(cfg); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	} else {
		defer database.CloseRedis()
	}

	// 初始化 JWT
	jwt.Init(cfg)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 设置路由
	routes.Setup(r)

	// 启动服务器
	addr := ":" + cfg.Server.Port
	log.Printf("Starting PhotoSet API server on %s", addr)
	log.Printf("Server mode: %s", cfg.Server.Mode)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

