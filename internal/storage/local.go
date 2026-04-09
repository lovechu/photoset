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
	os.MkdirAll(uploadDir, 0755)
	return &LocalStorage{UploadDir: uploadDir}
}

func (s *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	now := time.Now()
	dirPath := filepath.Join(s.UploadDir, "images", now.Format("2006"), now.Format("01"), now.Format("02"))
	os.MkdirAll(dirPath, 0755)

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

	return fmt.Sprintf("/uploads/images/%s/%s/%s/%s",
		now.Format("2006"), now.Format("01"), now.Format("02"), filename), nil
}
