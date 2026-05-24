package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostReportRepository handles database operations for post reports
type PostReportRepository struct {
	DB *gorm.DB
}

// NewPostReportRepository creates a new PostReportRepository
func NewPostReportRepository(db *gorm.DB) *PostReportRepository {
	return &PostReportRepository{DB: db}
}

// Create creates a new report
func (r *PostReportRepository) Create(report *domain.PostReport) error {
	return r.DB.Create(report).Error
}

// FindByID finds a report by ID
func (r *PostReportRepository) FindByID(id uint) (*domain.PostReport, error) {
	var report domain.PostReport
	err := r.DB.Preload("Post").Preload("Reply").Preload("Reporter").Preload("Handler").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// List returns all reports with pagination and status filter
func (r *PostReportRepository) List(page, pageSize int, status string) ([]domain.PostReport, int64, error) {
	var reports []domain.PostReport
	var total int64

	query := r.DB.Model(&domain.PostReport{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Post").Preload("Reply").Preload("Reporter").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// Update updates a report
func (r *PostReportRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.PostReport{}).Where("id = ?", id).Updates(updates).Error
}

// Process processes a report (resolve or reject)
func (r *PostReportRepository) Process(id uint, handlerID uint, status string, note string) error {
	updates := map[string]interface{}{
		"status":      status,
		"handler_id":  handlerID,
		"handled_at":  gorm.Expr("NOW()"),
		"handle_note": note,
	}
	return r.DB.Model(&domain.PostReport{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a report
func (r *PostReportRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.PostReport{}, id).Error
}
