package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Storage StorageConfig
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
			Secret:     getEnv("JWT_SECRET", "default-secret-key"),
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
			SignSecret:  getEnv("SIGN_SECRET", "default-sign-secret-change-me"),
			SignExpire:  getEnvAsInt("SIGN_EXPIRE", 7200),
		},
	}

	// ⚠️ 生产环境必须配置强密钥，默认值直接 panic
	if cfg.JWT.Secret == "default-secret-key" {
		log.Fatal("FATAL: JWT_SECRET is not configured. Set a strong random secret in .env or environment variable.")
	}
	if cfg.Storage.SignSecret == "default-sign-secret-change-me" {
		log.Fatal("FATAL: SIGN_SECRET is not configured. Set a strong random secret in .env or environment variable.")
	}

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
