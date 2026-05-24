package repository

import (
	"photoset/internal/domain"
	"time"

	"gorm.io/gorm"
)

// UserPointRepository handles database operations for user points
type UserPointRepository struct {
	DB *gorm.DB
}

// NewUserPointRepository creates a new UserPointRepository
func NewUserPointRepository(db *gorm.DB) *UserPointRepository {
	return &UserPointRepository{DB: db}
}

// GetByUserID gets user points by user ID
func (r *UserPointRepository) GetByUserID(userID uint) (*domain.UserPoint, error) {
	var userPoint domain.UserPoint
	err := r.DB.Where("user_id = ?", userID).First(&userPoint).Error
	if err != nil {
		// If not found, create a new one
		if err == gorm.ErrRecordNotFound {
			userPoint = domain.UserPoint{
				UserID: userID,
				Points: 0,
				Level:  1,
			}
			r.DB.Create(&userPoint)
			return &userPoint, nil
		}
		return nil, err
	}
	return &userPoint, nil
}

// Create creates a new user point record
func (r *UserPointRepository) Create(userPoint *domain.UserPoint) error {
	return r.DB.Create(userPoint).Error
}

// Update updates user points
func (r *UserPointRepository) Update(userPoint *domain.UserPoint) error {
	return r.DB.Save(userPoint).Error
}

// AddPoints adds points to a user (with transaction)
func (r *UserPointRepository) AddPoints(userID uint, points int) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var userPoint domain.UserPoint
		if err := tx.Where("user_id = ?", userID).First(&userPoint).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				userPoint = domain.UserPoint{
					UserID: userID,
					Points: 0,
					Level:  1,
				}
			} else {
				return err
			}
		}

		userPoint.Points += points
		if userPoint.Points < -9999 {
			userPoint.Points = -9999
		}
		userPoint.Level = domain.CalculateLevel(userPoint.Points)
		userPoint.UpdatedAt = time.Now()

		return tx.Save(&userPoint).Error
	})
}

// LogPointChange logs a point change
func (r *UserPointRepository) LogPointChange(userID uint, points int, action string, relatedID uint) error {
	log := map[string]interface{}{
		"user_id":    userID,
		"points":     points,
		"action":     action,
		"related_id":  relatedID,
		"created_at":  time.Now(), // Explicitly set timestamp
	}
	return r.DB.Table("user_point_logs").Create(&log).Error
}

// GetTodayPoints gets total points for a specific action today
func (r *UserPointRepository) GetTodayPoints(userID uint, action string) (int, error) {
	var total int
	
	// Use LIKE for SQLite date comparison (simpler and more compatible)
	today := time.Now().Format("2006-01-02")
	
	err := r.DB.Table("user_point_logs").
		Where("user_id = ? AND action = ? AND created_at LIKE ?", userID, action, today+"%").
		Select("COALESCE(SUM(points), 0)").
		Scan(&total).Error
	
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetTodayPostPoints gets points earned from posts today
func (r *UserPointRepository) GetTodayPostPoints(userID uint) (int, error) {
	return r.GetTodayPoints(userID, "post_create")
}

// GetTodayReplyPoints gets points earned from replies today
func (r *UserPointRepository) GetTodayReplyPoints(userID uint) (int, error) {
	return r.GetTodayPoints(userID, "reply_create")
}

// ListForAdmin returns all user points for admin with pagination and optional level filter
func (r *UserPointRepository) ListForAdmin(page, pageSize int, level int) ([]domain.UserPoint, int64, error) {
	var userPoints []domain.UserPoint
	var total int64

	query := r.DB.Model(&domain.UserPoint{})

	if level > 0 {
		query = query.Where("level = ?", level)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Order("points DESC").Offset(offset).Limit(pageSize).Find(&userPoints).Error
	return userPoints, total, err
}
