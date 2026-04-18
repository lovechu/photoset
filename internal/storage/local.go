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
	return s.UploadWithType(file, header, UploadTypePhoto, 0)
}

// UploadWithType 按类型分目录存储，photosetID 嵌入路径用于签名校验
// photo → /uploads/photos/{photosetID}/{MM}/{uuid}.ext
// cover → /uploads/covers/{photosetID}/{MM}/{uuid}.ext（photosetID=0 时用时间戳）
func (s *LocalStorage) UploadWithType(file multipart.File, header *multipart.FileHeader, ut UploadType, photosetID uint) (string, error) {
	subDir := "photos"
	if ut == UploadTypeCover {
		subDir = "covers"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	now := time.Now()

	// photosetID=0 时（封面先上传等场景），用时间戳做目录
	idOrDate := fmt.Sprintf("%d", photosetID)
	if photosetID == 0 {
		idOrDate = now.Format("20060102150405")
	}

	dirPath := filepath.Join(s.UploadDir, subDir, idOrDate, now.Format("01"))

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

	urlPath := fmt.Sprintf("/uploads/%s/%s/%s/%s",
		subDir, idOrDate, now.Format("01"), filename)

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
