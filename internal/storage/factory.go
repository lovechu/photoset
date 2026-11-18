package storage

import (
	"fmt"
	"net"
	"net/url"
	"photoset/internal/config"
	"strings"
)

// validateEndpoint 校验存储 endpoint，防止管理员配置被滥用于 SSRF。
// 仅允许 http/https 且指向公网地址，拒绝回环、内网、链路本地（含云元数据）地址。
func validateEndpoint(endpoint string) error {
	if endpoint == "" {
		return nil // R2 模式用 accountID 自动推导 endpoint
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("无效的 Endpoint 格式")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("Endpoint 仅支持 http/https 协议")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("Endpoint 缺少主机名")
	}
	// 拒绝 localhost 与常见内网/元数据域名
	lower := strings.ToLower(host)
	if lower == "localhost" || strings.HasSuffix(lower, ".localhost") ||
		lower == "metadata.google.internal" || lower == "169.254.169.254" {
		return fmt.Errorf("不允许使用本地或元数据地址作为 Endpoint")
	}
	// 解析为 IP 后判断是否为私有/回环/链路本地地址
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("无法解析 Endpoint 主机")
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("不允许使用内网或本地地址作为 Endpoint")
		}
	}
	return nil
}

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

	// SSRF 防护：拒绝内网/回环/元数据地址，避免连接测试被用作内网探测
	if err := validateEndpoint(endpoint); err != nil {
		return nil, err
	}

	return NewS3Storage(endpoint, region, accessKey, secretKey, bucket, publicURL, accountID)
}
