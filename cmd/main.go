package main

import (
	"crypto/rand"
	"encoding/hex"
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
	"gorm.io/gorm"

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
	// 注意：ensureSecrets 必须在 jwt.Init 之前执行，确保 cfg.JWT.Secret / cfg.Storage.SignSecret
	// 已通过三级回退（环境变量 > site_settings > 自动生成）赋值完成
	ensureSecrets(cfg)
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

// ensureSecrets 确保 JWT_SECRET 和 SIGN_SECRET 已配置，三级回退：
//  1. 环境变量（最高优先级，已有部署不受影响）
//  2. site_settings 表（重启后复用之前自动生成的值，多实例共享同一 DB 时密钥一致）
//  3. 首次部署：用 crypto/rand 生成 256bit 随机密钥，写入 site_settings 持久化
//
// 必须在 database.InitMySQL + AutoMigrate 之后、jwt.Init 之前调用。
// 这样既杜绝了使用公开已知弱默认值的风险，又让新部署完全零配置。
func ensureSecrets(cfg *config.Config) {
	db := database.GetMySQL()

	ensureOne := func(envKey, settingKey, label string, target *string) {
		// 第 1 级：环境变量已设置（config.Load 已读入），直接用
		if *target != "" {
			logger.Info("Secret loaded from environment variable", "key", envKey)
			return
		}
		// 第 2 级：从 site_settings 表加载（之前部署自动生成并持久化的值）
		if v, ok := getSetting(db, settingKey); ok {
			*target = v
			logger.Info("Secret loaded from site_settings", "key", settingKey)
			return
		}
		// 第 3 级：首次部署，自动生成并写入 DB
		secret := generateSecret()
		setSetting(db, settingKey, secret)
		*target = secret
		logger.Info("Secret auto-generated and persisted to site_settings", "key", settingKey)
	}

	ensureOne("JWT_SECRET", "jwt_secret", "JWT", &cfg.JWT.Secret)
	ensureOne("SIGN_SECRET", "sign_secret", "Sign", &cfg.Storage.SignSecret)

	// 最终兜底：理论上不会走到，防御性编程
	if cfg.JWT.Secret == "" {
		log.Fatal("FATAL: failed to initialize JWT_SECRET")
	}
	if cfg.Storage.SignSecret == "" {
		log.Fatal("FATAL: failed to initialize SIGN_SECRET")
	}
}

// generateSecret 使用 crypto/rand 生成 256bit 随机密钥（64 字符 hex）
// crypto/rand 走操作系统熵源，比 math/rand 更安全，适合生成密钥
func generateSecret() string {
	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("FATAL: failed to generate secret: %v", err)
	}
	return hex.EncodeToString(b)
}

// getSetting 从 site_settings 表按 key 读取单个值
// 返回 (value, found)；未找到或出错时 found=false
func getSetting(db *gorm.DB, key string) (string, bool) {
	var s domain.SiteSetting
	err := db.Where("`key` = ?", key).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", false
		}
		logger.Error("Failed to load secret from site_settings", "key", key, "error", err)
		return "", false
	}
	return s.Value, true
}

// setSetting 将密钥写入 site_settings 表（INSERT 或 UPDATE）
// 使用原生 SQL 的 ON DUPLICATE KEY UPDATE，避免 key 是 MySQL 保留字导致的语法错误
// 写入失败直接 fatal，因为这意味着后续 JWT 签发/验签会失败，无法安全运行
func setSetting(db *gorm.DB, key, value string) {
	sql := "INSERT INTO site_settings (`key`, `value`, `group`, created_at, updated_at) VALUES (?, ?, 'system', NOW(), NOW()) ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), updated_at = NOW()"
	if err := db.Exec(sql, key, value).Error; err != nil {
		log.Fatalf("FATAL: failed to persist secret to site_settings (key=%s): %v", key, err)
	}
}

