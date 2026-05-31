package main

import (
	"errors"
	"log"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/http/routes"
	"photoset/internal/logger"
	"photoset/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	// 初始化日志系统
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Info("Logger initialized", "level", cfg.Log.Level, "format", cfg.Log.Format)

	// 初始化 MySQL
	if err := database.InitMySQL(cfg); err != nil {
		logger.Error("Failed to initialize MySQL", "error", err)
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
		&domain.CommunityCategory{},
		&domain.Follow{},
		&domain.Notification{},
		&domain.Message{},
		&domain.PostShare{},
		&domain.Draft{},
		&domain.Topic{},
		// User level system
		&domain.UserLevelConfig{},
		&domain.Achievement{},
		&domain.UserAchievement{},
		&domain.PointsMallItem{},
		&domain.UserPointsExchange{},
		// User block/mute
		&domain.UserBlock{},
	); err != nil {
		// 忽略多对多关联表的重复主键错误（表已存在时 GORM 会尝试重复添加主键）
		if !isMultiplePrimaryKeyError(err) {
			log.Fatalf("Failed to auto migrate: %v", err)
		}
		logger.Warn("Migrate skipped duplicate primary key issue (safe to ignore)", "error", err)
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
			logger.Warn("Failed to create FULLTEXT index", "error", err)
		}
	}

	// 初始化 Redis（付费缓存依赖 Redis，必须成功）
	if err := database.InitRedis(cfg); err != nil {
		logger.Error("Failed to initialize Redis (required for paid cache)", "error", err)
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
	logger.Info("Starting PhotoSet API server", "addr", addr, "mode", cfg.Server.Mode)

	if err := r.Run(addr); err != nil {
		logger.Error("Failed to start server", "error", err)
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

