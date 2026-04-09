package storage

import "photoset/internal/config"

func NewStorage(cfg *config.StorageConfig) (Storage, error) {
	if cfg.Type == "r2" {
		return NewR2Storage(
			cfg.R2AccountID,
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			cfg.S3Bucket,
			cfg.R2PublicURL,
		)
	}
	return NewLocalStorage(cfg.LocalPath), nil
}
