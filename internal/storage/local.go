package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	UploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("⚠️ 无法创建上传目录 %s: %v\n", uploadDir, err)
	}
	return &LocalStorage{UploadDir: uploadDir}
}

func (s *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	now := time.Now()
	dirPath := filepath.Join(s.UploadDir, "images", now.Format("2006"), now.Format("01"), now.Format("02"))

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("无法创建存储目录: %v", err)
	}

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dirPath, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	urlPath := fmt.Sprintf("/uploads/images/%s/%s/%s/%s",
		now.Format("2006"), now.Format("01"), now.Format("02"), filename)

	return urlPath, nil
}

func (s *LocalStorage) TestConnection() error {
	// 测试本地存储目录是否可写
	testFile := filepath.Join(s.UploadDir, ".connection_test")
	if err := os.MkdirAll(s.UploadDir, 0755); err != nil {
		return fmt.Errorf("目录创建失败: %w", err)
	}
	if err := os.WriteFile(testFile, []byte("ok"), 0644); err != nil {
		return fmt.Errorf("目录不可写: %w", err)
	}
	os.Remove(testFile)
	return nil
}

func (s *LocalStorage) Type() StorageType {
	return StorageLocal
}
