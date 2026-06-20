package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"photoset/internal/config"
)

// validateBackupConfig 校验 DB 配置值，防止以 "-" 开头的可疑值被当作 mysqldump 选项
// （exec.Command 不走 shell，但仍是防御性编程，避免被工具误用）
func validateBackupConfig(cfg config.DBConfig) error {
	// 检查 Host / Port / User / Name 是否包含可疑的选项前缀
	fields := []struct {
		name  string
		value string
	}{
		{"host", cfg.Host},
		{"port", cfg.Port},
		{"user", cfg.User},
		{"name", cfg.Name},
	}
	for _, f := range fields {
		if strings.HasPrefix(f.value, "-") {
			return fmt.Errorf("数据库配置 %s 不能以 '-' 开头（防止参数注入）", f.name)
		}
	}
	return nil
}

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

	// 防御性校验：拒绝以 '-' 开头的配置值，防止被当作 mysqldump 选项
	if err := validateBackupConfig(s.cfg.DB); err != nil {
		return nil, fmt.Errorf("备份配置校验失败: %w", err)
	}

	// 构建 mysqldump 命令
	// 注意：不使用 -p 或 --password 参数传递密码，原因：
	//  1) -pPASSWORD 会在进程列表(ps)中暴露密码
	//  2) 空密码时 -p 会让 mysqldump 阻塞在 stdin 等待输入，导致服务挂起
	// 改用 MYSQL_PWD 环境变量（仅作用于子进程），mysqldump 5.6+ 原生支持
	cmd := exec.Command("mysqldump",
		"-h", s.cfg.DB.Host,
		"-P", s.cfg.DB.Port,
		"-u", s.cfg.DB.User,
		"--single-transaction",
		"--routines",
		"--triggers",
		s.cfg.DB.Name,
	)

	// 通过环境变量传递密码，避免出现在进程参数中
	// MYSQL_PWD 仅对当前子进程及其后代生效，不会污染主进程环境
	cmd.Env = append(os.Environ(), "MYSQL_PWD="+s.cfg.DB.Password)

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
