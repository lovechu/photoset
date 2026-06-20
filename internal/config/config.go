package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Storage    StorageConfig
	Log        LogConfig
	Alipay     AlipayConfig
	WechatPay  WechatPayConfig
}

type LogConfig struct {
	Level  string // debug/info/warn/error
	Format string // json/text
}

type ServerConfig struct {
	Port string
	Mode string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Charset  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type StorageConfig struct {
	Type        string
	LocalPath   string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
	R2AccountID string // Cloudflare Account ID
	R2PublicURL string // R2 自定义域名，如 https://assets.example.com
	SignSecret  string // URL 签名密钥
	SignExpire  int    // 签名过期时间（秒），默认 7200（2小时）
}

type AlipayConfig struct {
	AppID      string // 支付宝应用ID
	PrivateKey string // 应用私钥
	PublicKey  string // 支付宝公钥
	NotifyURL  string // 异步通知地址
	ReturnURL  string // 同步跳转地址
	IsSandbox  bool   // 是否沙箱环境
}

type WechatPayConfig struct {
	AppID    string // 微信开放平台 AppID
	MchID    string // 微信支付商户号
	APIKey   string // API 密钥 (APIv2)
	CertPath string // 证书路径
	NotifyURL string // 异步通知地址
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "photoset"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			// 默认空字符串：由 main.go 的 ensureSecrets 在 DB 初始化后做三级回退
			// （环境变量 > site_settings > 自动生成），避免使用公开已知的弱默认值
			Secret:     getEnv("JWT_SECRET", ""),
			ExpireHour: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Storage: StorageConfig{
			Type:        getEnv("STORAGE_TYPE", "local"),
			LocalPath:   getEnv("LOCAL_STORAGE_PATH", "./uploads"),
			S3Endpoint:  getEnv("S3_ENDPOINT", ""),
			S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey: getEnv("S3_SECRET_KEY", ""),
			S3Bucket:    getEnv("S3_BUCKET", ""),
			S3Region:    getEnv("S3_REGION", ""),
			R2AccountID: getEnv("R2_ACCOUNT_ID", ""),
			R2PublicURL: getEnv("R2_PUBLIC_URL", ""),
			// 默认空字符串：由 main.go 的 ensureSecrets 在 DB 初始化后做三级回退
			SignSecret:  getEnv("SIGN_SECRET", ""),
			SignExpire:  getEnvAsInt("SIGN_EXPIRE", 7200),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
		Alipay: AlipayConfig{
			AppID:      getEnv("ALIPAY_APP_ID", ""),
			PrivateKey: getEnv("ALIPAY_PRIVATE_KEY", ""),
			PublicKey:  getEnv("ALIPAY_PUBLIC_KEY", ""),
			NotifyURL:  getEnv("ALIPAY_NOTIFY_URL", ""),
			ReturnURL:  getEnv("ALIPAY_RETURN_URL", ""),
			IsSandbox:  getEnv("ALIPAY_SANDBOX", "false") == "true",
		},
		WechatPay: WechatPayConfig{
			AppID:     getEnv("WECHAT_APP_ID", ""),
			MchID:     getEnv("WECHAT_MCH_ID", ""),
			APIKey:    getEnv("WECHAT_API_KEY", ""),
			CertPath:  getEnv("WECHAT_CERT_PATH", ""),
			NotifyURL: getEnv("WECHAT_NOTIFY_URL", ""),
		},
	}

	// 密钥校验已移至 main.go 的 ensureSecrets：
	// 在 DB 初始化后做三级回退（环境变量 > site_settings > 自动生成），
	// 最终仍为空才 panic，避免 config.Load 阶段无法访问 DB 的限制。

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
