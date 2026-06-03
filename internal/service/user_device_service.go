package service

import (
	"time"

	"photoset/internal/domain"
	"photoset/internal/repository"
)

// UserDeviceService handles user device business logic
type UserDeviceService struct {
	deviceRepo *repository.UserDeviceRepository
}

// NewUserDeviceService creates a new UserDeviceService
func NewUserDeviceService(deviceRepo *repository.UserDeviceRepository) *UserDeviceService {
	return &UserDeviceService{deviceRepo: deviceRepo}
}

// RegisterOrUpdateDevice registers a new device or updates existing one
// deviceName is the user-friendly name from the frontend (e.g. "Flutter Mobile")
// userAgent is the raw User-Agent header used to parse device type, OS, and browser
func (s *UserDeviceService) RegisterOrUpdateDevice(userID uint, deviceID, deviceName, userAgent, ip, ipLocation string) error {
	// Parse user agent to extract device info
	parsedName, deviceType, os, browser := parseDeviceInfo(userAgent)

	// Use frontend-provided device name if available, otherwise use parsed name
	if deviceName == "" {
		deviceName = parsedName
	}

	// Check if device already exists
	existing, err := s.deviceRepo.FindByDeviceID(userID, deviceID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Update existing device
		existing.LastActiveAt = time.Now()
		existing.IP = ip
		existing.IPLocation = ipLocation
		existing.UserAgent = userAgent
		existing.OS = os
		existing.Browser = browser
		existing.IsActive = true
		existing.DeviceName = deviceName
		return s.deviceRepo.Save(existing)
	}

	// Create new device
	device := &domain.UserDevice{
		UserID:       userID,
		DeviceID:     deviceID,
		DeviceName:   deviceName,
		DeviceType:   deviceType,
		OS:           os,
		Browser:      browser,
		IP:           ip,
		IPLocation:   ipLocation,
		UserAgent:    userAgent,
		LastActiveAt: time.Now(),
		IsActive:     true,
	}

	return s.deviceRepo.Create(device)
}

// GetUserDevices gets all active devices for a user
func (s *UserDeviceService) GetUserDevices(userID uint) ([]domain.UserDevice, error) {
	return s.deviceRepo.GetByUserID(userID)
}

// DeactivateDevice deactivates a specific device
func (s *UserDeviceService) DeactivateDevice(userID uint, deviceID string) error {
	return s.deviceRepo.DeactivateDevice(userID, deviceID)
}

// DeactivateOtherDevices deactivates all devices except the current one
func (s *UserDeviceService) DeactivateOtherDevices(userID uint, currentDeviceID string) error {
	return s.deviceRepo.DeactivateOtherDevices(userID, currentDeviceID)
}

// CleanupInactiveDevices deletes devices inactive for more than specified days
func (s *UserDeviceService) CleanupInactiveDevices(days int) error {
	return s.deviceRepo.DeleteOldDevices(days)
}

// parseDeviceInfo extracts device name, type, OS, and browser from user agent
func parseDeviceInfo(userAgent string) (deviceName, deviceType, os, browser string) {
	// Reuse the parseUserAgent function from login_history_service
	// but we need to import it or duplicate the logic
	// For simplicity, we'll call the same parsing logic
	device, browserResult, osResult := parseUserAgent(userAgent)
	
	// Generate device name based on device type and OS
	deviceName = device
	if osResult != "Unknown" && osResult != "" {
		deviceName = device + " (" + osResult + ")"
	}
	
	deviceType = device
	
	return deviceName, deviceType, osResult, browserResult
}