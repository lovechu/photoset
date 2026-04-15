package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type StorageType string

const (
	StorageLocal StorageType = "local"
	StorageS3    StorageType = "s3" // 通用 S3 兼容存储（含 Cloudflare R2、阿里 OSS、MinIO 等）
)

type Storage interface {
	Upload(file multipart.File, header *multipart.FileHeader) (string, error)
	TestConnection() error
	Type() StorageType
}

// S3Storage 通用 S3 兼容对象存储（支持 Cloudflare R2、阿里 OSS、MinIO 等）
type S3Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string // 自定义域名或 CDN 域名，如 https://assets.example.com
	storageType StorageType
}

// NewS3Storage 创建通用 S3 存储
// R2 特殊处理：如果 endpoint 为空且 accountID 不为空，自动生成 R2 endpoint
func NewS3Storage(endpoint, region, accessKey, secretKey, bucket, publicURL, accountID string) (*S3Storage, error) {
	// R2 兼容：如果没有 endpoint 但有 accountID，自动生成 R2 endpoint
	if endpoint == "" && accountID != "" {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
		region = "auto"
	}
	if region == "" {
		region = "us-east-1"
	}

	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(region))
	opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")))

	if endpoint != "" {
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: endpoint}, nil
		})
		opts = append(opts, config.WithEndpointResolverWithOptions(resolver))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, fmt.Errorf("S3 配置加载失败: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{
		client:      client,
		bucket:      bucket,
		publicURL:   publicURL,
		storageType: StorageS3,
	}, nil
}

func (s *S3Storage) Upload(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	now := time.Now()
	key := fmt.Sprintf("images/%s/%s/%s/%s%s",
		now.Format("2006"), now.Format("01"), now.Format("02"),
		uuid.New().String(), ext)

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mimeType(ext)),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(s.publicURL, "/"), key), nil
}

// TestConnection 测试存储连接
func (s *S3Storage) TestConnection() error {
	_, err := s.client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("连接测试失败: %w", err)
	}
	return nil
}

func (s *S3Storage) Type() StorageType {
	return s.storageType
}

func mimeType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
