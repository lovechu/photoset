package storage

import (
	"fmt"
	"photoset/internal/config"
)

func NewStorage(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "s3", "r2":
		return NewS3Storage(
			cfg.S3Endpoint,
			cfg.S3Region,
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			cfg.S3Bucket,
			cfg.R2PublicURL,
			cfg.R2AccountID,
		)
	default:
		return NewLocalStorage(cfg.LocalPath), nil
	}
}

// NewStorageFromSettings 从站点设置动态创建存储实例（用于测试连接）
func NewStorageFromSettings(settings map[string]interface{}) (Storage, error) {
	storageType, _ := settings["storage_type"].(string)
	if storageType != "s3" && storageType != "r2" {
		return NewLocalStorage("./uploads"), nil
	}

	endpoint, _ := settings["s3_endpoint"].(string)
	region, _ := settings["s3_region"].(string)
	accessKey, _ := settings["s3_access_key"].(string)
	secretKey, _ := settings["s3_secret_key"].(string)
	bucket, _ := settings["s3_bucket"].(string)
	publicURL, _ := settings["cdn_domain"].(string)
	accountID, _ := settings["r2_account_id"].(string)

	if endpoint == "" && accountID == "" {
		return nil, fmt.Errorf("请填写 S3 Endpoint 或 R2 Account ID")
	}
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("请填写 Access Key 和 Secret Key")
	}
	if bucket == "" {
		return nil, fmt.Errorf("请填写 Bucket 名称")
	}

	return NewS3Storage(endpoint, region, accessKey, secretKey, bucket, publicURL, accountID)
}
