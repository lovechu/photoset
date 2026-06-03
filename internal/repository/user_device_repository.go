package repository

import (
	"time"

	"photoset/internal/domain"

	"gorm.io/gorm"
)

// UserDeviceRepository handles user device data operations
type UserDeviceRepository struct {
	db *gorm.DB
}

// NewUserDeviceRepository creates a new UserDeviceRepository
func NewUserDeviceRepository(db *gorm.DB) *UserDeviceRepository {
	return &UserDeviceRepository{db: db}
}

// Create adds a new device record
func (r *UserDeviceRepository) Create(device *domain.UserDevice) error {
	return r.db.Create(device).Error
}

// GetByUserID gets devices for a user
func (r *UserDeviceRepository) GetByUserID(userID uint) ([]domain.UserDevice, error) {
	var devices []domain.UserDevice
	if err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_active_at DESC").
		Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// FindByDeviceID finds a device by user ID and device ID
func (r *UserDeviceRepository) FindByDeviceID(userID uint, deviceID string) (*domain.UserDevice, error) {
	var device domain.UserDevice
	err := r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).
		First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// UpdateLastActive updates the last active time
func (r *UserDeviceRepository) UpdateLastActive(userID uint, deviceID string) error {
	return r.db.Model(&domain.UserDevice{}).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Update("last_active_at", time.Now()).Error
}

// DeactivateDevice deactivates a device
func (r *UserDeviceRepository) DeactivateDevice(userID uint, deviceID string) error {
	return r.db.Model(&domain.UserDevice{}).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Update("is_active", false).Error
}

// DeactivateOtherDevices deactivates all other devices except the current one
func (r *UserDeviceRepository) DeactivateOtherDevices(userID uint, currentDeviceID string) error {
	return r.db.Model(&domain.UserDevice{}).
		Where("user_id = ? AND device_id != ?", userID, currentDeviceID).
		Update("is_active", false).Error
}

// DeleteOldDevices deletes devices inactive for more than days
func (r *UserDeviceRepository) DeleteOldDevices(days int) error {
	return r.db.Where("last_active_at < NOW() - INTERVAL ? DAY AND is_active = ?", days, false).
		Delete(&domain.UserDevice{}).Error
}
