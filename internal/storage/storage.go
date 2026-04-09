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
	StorageR2    StorageType = "r2"
)

type Storage interface {
	Upload(file multipart.File, header *multipart.FileHeader) (string, error)
}

// R2Storage Cloudflare R2 存储
type R2Storage struct {
	client    *s3.Client
	bucket    string
	publicURL string // R2 自定义域名，如 https://assets.example.com
}

func NewR2Storage(accountID, accessKey, secretKey, bucket, publicURL string) (*R2Storage, error) {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &R2Storage{client: client, bucket: bucket, publicURL: publicURL}, nil
}

func (s *R2Storage) Upload(file multipart.File, header *multipart.FileHeader) (string, error) {
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

func mimeType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
