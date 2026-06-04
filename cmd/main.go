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

	_ "photoset/docs" // Swagger docs
)

// @title           PhotoSet API
// @version         1.0
// @description     PhotoSet 摄影套图浏览平台后端 API 文档
// @description     提供套图管理、社区互动、用户认证、支付等功能
// @termsOfService  https://photoset.example.com/terms

// @contact.name   API Support
// @contact.email  support@photoset.example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在请求头中添加 Bearer {token} 进行认证

// @tag.name Auth
// @tag.description 用户认证相关接口（注册、登录、密码重置等）

// @tag.name PhotoSet
// @tag.description 套图管理接口（CRUD、列表、详情、下载等）

// @tag.name Community
// @tag.description 社区互动接口（帖子、回帖、点赞、收藏等）

// @tag.name User
// @tag.description 用户信息接口

// @tag.name Favorites
// @tag.description 收藏夹接口

// @tag.name Comments
// @tag.description 评论接口

// @tag.name Collections
// @tag.description 合集接口

// @tag.name Orders
// @tag.description 订单与支付接口

// @tag.name Membership
// @tag.description 会员套餐接口

// @tag.name Upload
// @tag.description 文件上传接口

// @tag.name Tags
// @tag.description 标签管理接口

// @tag.name Categories
// @tag.description 分类管理接口

// @tag.name Explore
// @tag.description 探索发现接口

// @tag.name Creator
// @tag.description 创作者数据统计接口

// @tag.name Feedback
// @tag.description 用户反馈接口

// @tag.name ViewHistory
// @tag.description 浏览历史接口

// @tag.name Notification
// @tag.description 通知接口

// @tag.name Message
// @tag.description 私信接口

// @tag.name Follow
// @tag.description 关注接口

// @tag.name UserLevel
// @tag.description 用户等级与积分接口

// @tag.name OAuth2
// @tag.description OAuth2 第三方授权接口

// @tag.name Share
// @tag.description 分享链接接口

// @tag.name Review
// @tag.description 套图评价接口

// @tag.name Pages
// @tag.description 页面管理接口

// @tag.name Settings
// @tag.description 站点设置接口

// @tag.name System
// @tag.description 系统管理接口

// @tag.name Admin
// @tag.description 管理后台接口

// @tag.name Health
// @tag.description 健康检查接口

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
		// OAuth2 third-party authorization
		&domain.OAuthClient{},
		&domain.OAuthAuthorization{},
		&domain.OAuthToken{},
		// User security & privacy
		&domain.LoginHistory{},
		&domain.UserDevice{},
		&domain.UserPrivacySetting{},
		// Browsing history
		&domain.ViewHistory{},
		// Email verification
		&domain.EmailVerificationCode{},
		// User collections
		&domain.UserCollection{},
		&domain.CollectionItem{},
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

