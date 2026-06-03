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
	fmt.Printf("[STORAGE-DEBUG] LocalStorage.UploadWithType 开始: type=%s, photosetID=%d\n", ut, photosetID)

	subDir := "photos"
	if ut == UploadTypeCover {
		subDir = "covers"
	} else if ut == UploadTypeVideo {
		subDir = "videos"
	} else if ut == UploadTypeAvatar {
		subDir = "avatars"
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	now := time.Now()

	idOrDate := fmt.Sprintf("%d", photosetID)
	if photosetID == 0 {
		idOrDate = now.Format("20060102150405")
	}

	dirPath := filepath.Join(s.UploadDir, subDir, idOrDate, now.Format("01"))
	fmt.Printf("[STORAGE-DEBUG] 目标目录: %s\n", dirPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Printf("[STORAGE-DEBUG] 创建目录失败: %v\n", err)
		return "", fmt.Errorf("无法创建存储目录: %v", err)
	}
	fmt.Printf("[STORAGE-DEBUG] 目录创建成功\n")

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dirPath, filename)
	fmt.Printf("[STORAGE-DEBUG] 文件路径: %s\n", fullPath)

	dst, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("[STORAGE-DEBUG] 创建文件失败: %v\n", err)
		return "", err
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		fmt.Printf("[STORAGE-DEBUG] 写入文件失败: %v\n", err)
		return "", err
	}
	fmt.Printf("[STORAGE-DEBUG] 写入成功: %d bytes\n", written)

	urlPath := fmt.Sprintf("/uploads/%s/%s/%s/%s",
		subDir, idOrDate, now.Format("01"), filename)
	fmt.Printf("[STORAGE-DEBUG] 上传完成: %s\n", urlPath)

	return urlPath, nil
}

// Delete 根据存储 URL 删除本地物理文件
// url 格式：/uploads/photos/123/uuid.jpg → 转换为本地路径
func (s *LocalStorage) Delete(url string) error {
	if url == "" {
		return nil
	}
	// 去掉 URL 中的 /uploads/ 前缀，拼接为本地路径
	relPath := url
	relPath = strings.TrimPrefix(relPath, "/uploads/")
	relPath = strings.TrimPrefix(relPath, "/uploads")
	fullPath := filepath.Join(s.UploadDir, relPath)

	// 安全检查：确保路径在 UploadDir 内
	absUploadDir, _ := filepath.Abs(s.UploadDir)
	absFullPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFullPath, absUploadDir) {
		return fmt.Errorf("路径越界: %s", url)
	}

	if err := os.Remove(absFullPath); err != nil {
		if os.IsNotExist(err) {
			return nil // 文件已不存在，视为成功
		}
		return fmt.Errorf("删除文件失败 [%s]: %w", absFullPath, err)
	}
	return nil
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
