package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/logger"
	"photoset/internal/pkg/password"
)

const (
	// DeletionCooldownDays 冷静期天数
	DeletionCooldownDays = 30
)

type AccountDeletionService struct {
	db *gorm.DB
}

func NewAccountDeletionService() *AccountDeletionService {
	return &AccountDeletionService{
		db: database.GetMySQL(),
	}
}

// RequestDeletion 申请注销账号
// 设置 deletion_requested_at，30天后定时任务将执行真正的删除
func (s *AccountDeletionService) RequestDeletion(userID uint, password string, reason string) error {
	var user domain.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证密码
	if !checkPassword(password, user.PasswordHash) {
		return errors.New("密码错误")
	}

	// 检查是否已申请注销
	if user.DeletionRequestedAt != nil {
		return errors.New("您已申请注销，请等待处理")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"deletion_requested_at": now,
		"deletion_reason":       reason,
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		logger.Error("申请注销失败", "error", err, "user_id", userID)
		return errors.New("申请注销失败，请稍后重试")
	}

	logger.Info("用户申请注销", "user_id", userID, "reason", reason)
	return nil
}

// CancelDeletion 取消注销申请
func (s *AccountDeletionService) CancelDeletion(userID uint) error {
	var user domain.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.DeletionRequestedAt == nil {
		return errors.New("您尚未申请注销")
	}

	updates := map[string]interface{}{
		"deletion_requested_at": nil,
		"deletion_reason":       "",
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		logger.Error("取消注销失败", "error", err, "user_id", userID)
		return errors.New("取消注销失败，请稍后重试")
	}

	logger.Info("用户取消注销", "user_id", userID)
	return nil
}

// GetDeletionStatus 获取注销状态
func (s *AccountDeletionService) GetDeletionStatus(userID uint) (map[string]interface{}, error) {
	var user domain.User
	if err := s.db.Select("id", "deletion_requested_at").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	result := map[string]interface{}{
		"has_pending_deletion": user.DeletionRequestedAt != nil,
	}

	if user.DeletionRequestedAt != nil {
		deadline := user.DeletionRequestedAt.AddDate(0, 0, DeletionCooldownDays)
		result["requested_at"] = user.DeletionRequestedAt.Format(time.RFC3339)
		result["deadline"] = deadline.Format(time.RFC3339)
		result["days_remaining"] = int(time.Until(deadline).Hours()/24) + 1
	}

	return result, nil
}

// ProcessExpiredDeletions 处理过期的注销申请（定时任务调用）
// 删除冷静期已过的用户数据
func (s *AccountDeletionService) ProcessExpiredDeletions() (int, error) {
	deadline := time.Now().AddDate(0, 0, -DeletionCooldownDays)

	var users []domain.User
	if err := s.db.Where("deletion_requested_at IS NOT NULL AND deletion_requested_at <= ?", deadline).
		Find(&users).Error; err != nil {
		return 0, err
	}

	processed := 0
	for _, user := range users {
		// 使用事务确保数据一致性
		tx := s.db.Begin()

		// 1. 匿名化用户数据
		updates := map[string]interface{}{
			"nickname":              "已注销用户",
			"email":                 fmt.Sprintf("deleted_%d@deleted.com", user.ID),
			"password_hash":         "",
			"avatar":                "",
			"bio":                   "",
			"status":                0, // 禁用
			"deletion_requested_at": nil,
			"deletion_reason":       "",
		}

		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			tx.Rollback()
			logger.Error("匿名化用户数据失败", "error", err, "user_id", user.ID)
			continue
		}

		// 2. 软删除用户
		if err := tx.Delete(&user).Error; err != nil {
			tx.Rollback()
			logger.Error("软删除用户失败", "error", err, "user_id", user.ID)
			continue
		}

		// 3. 删除相关数据（浏览历史、收藏等）
		tx.Where("user_id = ?", user.ID).Delete(&domain.ViewHistory{})

		tx.Commit()
		processed++
		logger.Info("已处理过期注销用户", "user_id", user.ID)
	}

	return processed, nil
}

// checkPassword 验证密码
func checkPassword(pwd, hash string) bool {
	return password.Check(pwd, hash)
}