package main

import (
	"errors"
	"log"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/http/routes"
	"photoset/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	// 初始化 MySQL
	if err := database.InitMySQL(cfg); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}
	defer database.CloseMySQL()

	// Auto migrate (migrate non-associated tables first)
	if err := database.GetMySQL().AutoMigrate(
		&domain.User{},
		&domain.PhotoSet{},
		&domain.Photo{},
		&domain.Tag{},
		&domain.Category{},
		&domain.Favorite{},
		&domain.MembershipPlan{},
		&domain.Order{},
		&domain.SiteSetting{},
		&domain.Page{},
		&domain.AdminLog{},
		&domain.PasswordResetToken{},
		&domain.ApiKey{},
		&domain.Comment{},
		&domain.CommentLike{},
		// Community module
		&domain.Post{},
		&domain.PostReply{},
		&domain.PostLike{},
		&domain.PostReplyLike{},
		&domain.UserPoint{},
		&domain.SensitiveWord{},
		&domain.PostReport{},
	); err != nil {
		// 忽略多对多关联表的重复主键错误（表已存在时 GORM 会尝试重复添加主键）
		if !isMultiplePrimaryKeyError(err) {
			log.Fatalf("Failed to auto migrate: %v", err)
		}
		log.Printf("Warning: migrate skipped duplicate primary key issue (safe to ignore): %v", err)
	}

	// 确保 FULLTEXT 索引存在（容错方式）
	var count int64
	database.GetMySQL().Raw(`
		SELECT COUNT(*) FROM information_schema.STATISTICS
		WHERE table_schema = DATABASE() AND table_name = 'photosets'
		AND index_name = 'ft_title_description'
	`).Scan(&count)
	if count == 0 {
		if err := database.GetMySQL().Exec(`
			CREATE FULLTEXT INDEX ft_title_description
			ON photosets (title, description) WITH PARSER ngram
		`).Error; err != nil {
			log.Printf("Warning: Failed to create FULLTEXT index: %v", err)
		}
	}

	// 初始化 Redis（付费缓存依赖 Redis，必须成功）
	if err := database.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to initialize Redis (required for paid cache): %v", err)
	}
	defer database.CloseRedis()

	// 初始化 JWT
	jwt.Init(cfg)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 设置路由
	routes.Setup(r, cfg)

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

// isMultiplePrimaryKeyError 检测是否是重复主键的 MySQL 错误（Error #1068）
func isMultiplePrimaryKeyError(err error) bool {
	if err == nil {
		return false
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1068 // ER_DUP_KEYNAME
	}
	return false
}

