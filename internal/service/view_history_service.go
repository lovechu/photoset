package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
	"time"
)

type ViewHistoryService struct {
	repo *repository.ViewHistoryRepository
}

func NewViewHistoryService(repo *repository.ViewHistoryRepository) *ViewHistoryService {
	return &ViewHistoryService{repo: repo}
}

// Record 记录浏览历史
func (s *ViewHistoryService) Record(userID, photosetID uint, title, coverImage string, photoCount int, isVip bool, price float64) error {
	history := &domain.ViewHistory{
		UserID:     userID,
		PhotosetID: photosetID,
		ViewedAt:   time.Now(),
		Title:      title,
		CoverImage: coverImage,
		PhotoCount: photoCount,
		IsVip:      isVip,
		Price:      price,
	}
	return s.repo.Record(history)
}

// List 分页查询浏览历史
func (s *ViewHistoryService) List(userID uint, page, pageSize int) ([]domain.ViewHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.List(userID, page, pageSize)
}

// ListByDateRange 按日期范围查询浏览历史
func (s *ViewHistoryService) ListByDateRange(userID uint, start, end time.Time, page, pageSize int) ([]domain.ViewHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.ListByDateRange(userID, start, end, page, pageSize)
}

// Delete 删除单条浏览历史
func (s *ViewHistoryService) Delete(userID, historyID uint) error {
	return s.repo.Delete(userID, historyID)
}

// ClearAll 清空用户所有浏览历史
func (s *ViewHistoryService) ClearAll(userID uint) error {
	return s.repo.ClearAll(userID)
}

// BatchDelete 批量删除浏览历史
func (s *ViewHistoryService) BatchDelete(userID uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BatchDelete(userID, ids)
}
