package repository

import (
	"photoset/internal/domain"
	"time"

	"gorm.io/gorm"
)

type ViewHistoryRepository struct {
	db *gorm.DB
}

func NewViewHistoryRepository(db *gorm.DB) *ViewHistoryRepository {
	return &ViewHistoryRepository{db: db}
}

// Record 记录浏览历史（去重：同一用户同一套图只保留最新一条）
func (r *ViewHistoryRepository) Record(history *domain.ViewHistory) error {
	// 先删除已有的记录
	r.db.Where("user_id = ? AND photoset_id = ?", history.UserID, history.PhotosetID).Delete(&domain.ViewHistory{})
	// 创建新记录
	return r.db.Create(history).Error
}

// List 分页查询浏览历史（按浏览时间倒序）
func (r *ViewHistoryRepository) List(userID uint, page, pageSize int) ([]domain.ViewHistory, int64, error) {
	var histories []domain.ViewHistory
	var total int64

	query := r.db.Model(&domain.ViewHistory{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("viewed_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&histories).Error

	return histories, total, err
}

// ListByDateRange 按日期范围查询浏览历史
func (r *ViewHistoryRepository) ListByDateRange(userID uint, start, end time.Time, page, pageSize int) ([]domain.ViewHistory, int64, error) {
	var histories []domain.ViewHistory
	var total int64

	query := r.db.Model(&domain.ViewHistory{}).
		Where("user_id = ? AND viewed_at >= ? AND viewed_at < ?", userID, start, end)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("viewed_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&histories).Error

	return histories, total, err
}

// Delete 删除单条浏览历史
func (r *ViewHistoryRepository) Delete(userID, historyID uint) error {
	return r.db.Where("user_id = ? AND id = ?", userID, historyID).Delete(&domain.ViewHistory{}).Error
}

// ClearAll 清空用户所有浏览历史
func (r *ViewHistoryRepository) ClearAll(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.ViewHistory{}).Error
}

// BatchDelete 批量删除浏览历史
func (r *ViewHistoryRepository) BatchDelete(userID uint, ids []uint) error {
	return r.db.Where("user_id = ? AND id IN ?", userID, ids).Delete(&domain.ViewHistory{}).Error
}
