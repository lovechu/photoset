package repository

import (
	"photoset/internal/database"
	"photoset/internal/domain"
)

type AdminLogRepository struct{}

func NewAdminLogRepository() *AdminLogRepository {
	return &AdminLogRepository{}
}

// Create 记录一条操作日志
func (r *AdminLogRepository) Create(log *domain.AdminLog) error {
	return database.GetMySQL().Create(log).Error
}

// List 分页查询操作日志
func (r *AdminLogRepository) List(page, pageSize int, action string) ([]domain.AdminLog, int64, error) {
	db := database.GetMySQL().Model(&domain.AdminLog{})
	if action != "" {
		db = db.Where("action = ?", action)
	}
	var total int64
	db.Count(&total)

	var logs []domain.AdminLog
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}
