package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"photoset/internal/config"
)

type BackupService struct {
	cfg       *config.Config
	backupDir string
}

type BackupInfo struct {
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	SizeStr   string    `json:"size_str"`
	CreatedAt time.Time `json:"created_at"`
}

func NewBackupService(cfg *config.Config) *BackupService {
	backupDir := "./backups"
	os.MkdirAll(backupDir, 0755)

	return &BackupService{
		cfg:       cfg,
		backupDir: backupDir,
	}
}

// CreateBackup 创建数据库备份
func (s *BackupService) CreateBackup() (*BackupInfo, error) {
	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("photoset_%s.sql", timestamp)
	filepath := filepath.Join(s.backupDir, filename)

	// 构建 mysqldump 命令
	cmd := exec.Command("mysqldump",
		"-h", s.cfg.DB.Host,
		"-P", s.cfg.DB.Port,
		"-u", s.cfg.DB.User,
		"-p"+s.cfg.DB.Password,
		"--single-transaction",
		"--routines",
		"--triggers",
		s.cfg.DB.Name,
	)

	// 创建输出文件
	outFile, err := os.Create(filepath)
	if err != nil {
		return nil, fmt.Errorf("创建备份文件失败: %w", err)
	}
	defer outFile.Close()

	// 设置命令输出到文件
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	// 执行备份命令
	if err := cmd.Run(); err != nil {
		os.Remove(filepath) // 清理失败的备份文件
		return nil, fmt.Errorf("执行 mysqldump 失败: %w", err)
	}

	// 获取文件信息
 fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("获取备份文件信息失败: %w", err)
	}

	return &BackupInfo{
		Filename:  filename,
		Size:      fileInfo.Size(),
		SizeStr:   formatFileSize(fileInfo.Size()),
		CreatedAt: fileInfo.ModTime(),
	}, nil
}

// ListBackups 获取备份列表
func (s *BackupService) ListBackups() ([]BackupInfo, error) {
	entries, err := os.ReadDir(s.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupInfo{}, nil
		}
		return nil, fmt.Errorf("读取备份目录失败: %w", err)
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, BackupInfo{
			Filename:  entry.Name(),
			Size:      info.Size(),
			SizeStr:   formatFileSize(info.Size()),
			CreatedAt: info.ModTime(),
		})
	}

	// 按创建时间倒序排列
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

// GetBackupPath 获取备份文件路径
func (s *BackupService) GetBackupPath(filename string) (string, error) {
	// 防止路径遍历攻击
	if filepath.Base(filename) != filename {
		return "", fmt.Errorf("无效的文件名")
	}

	path := filepath.Join(s.backupDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("备份文件不存在")
	}

	return path, nil
}

// DeleteBackup 删除备份文件
func (s *BackupService) DeleteBackup(filename string) error {
	// 防止路径遍历攻击
	if filepath.Base(filename) != filename {
		return fmt.Errorf("无效的文件名")
	}

	path := filepath.Join(s.backupDir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在")
	}

	return os.Remove(path)
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
