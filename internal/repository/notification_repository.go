package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	DB *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// Create creates a new notification
func (r *NotificationRepository) Create(notification *domain.Notification) error {
	return r.DB.Create(notification).Error
}

// FindByID finds a notification by ID
func (r *NotificationRepository) FindByID(id uint) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.DB.Preload("Sender").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ListByUserID returns notifications for a user with pagination
func (r *NotificationRepository) ListByUserID(userID uint, page, pageSize int, unreadOnly bool) ([]domain.Notification, int64, error) {
	var notifications []domain.Notification
	var total int64

	query := r.DB.Model(&domain.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Sender").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// CountUnread counts unread notifications for a user
func (r *NotificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(id uint, userID uint) error {
	return r.DB.Model(&domain.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.DB.Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

// Delete deletes a notification (soft delete)
func (r *NotificationRepository) Delete(id uint, userID uint) error {
	return r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Notification{}).Error
}

// DeleteAll deletes all notifications for a user
func (r *NotificationRepository) DeleteAll(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&domain.Notification{}).Error
}
