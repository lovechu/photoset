package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
	"strings"
)

// LoginHistoryService handles login history business logic
type LoginHistoryService struct {
	historyRepo *repository.LoginHistoryRepository
}

// NewLoginHistoryService creates a new LoginHistoryService
func NewLoginHistoryService(historyRepo *repository.LoginHistoryRepository) *LoginHistoryService {
	return &LoginHistoryService{historyRepo: historyRepo}
}

// CreateLoginHistory records a login attempt
func (s *LoginHistoryService) CreateLoginHistory(userID uint, ip, ipLocation, userAgent, loginType, deviceType, osVersion string, success bool, failReason string) error {
	// Parse user agent to extract device info
	uaDevice, browser, uaOS := parseUserAgent(userAgent)

	// Prioritize frontend-provided device info over UA parsing
	device := uaDevice
	os := uaOS
	if deviceType != "" {
		device = deviceType
	}
	if osVersion != "" {
		os = osVersion
	}

	history := &domain.LoginHistory{
		UserID:     userID,
		IP:         ip,
		IPLocation: ipLocation,
		UserAgent:  userAgent,
		Device:     device,
		Browser:    browser,
		OS:         os,
		LoginType:  loginType,
		Success:    success,
		FailReason: failReason,
	}

	return s.historyRepo.Create(history)
}

// GetLoginHistory gets login history for a user with pagination
func (s *LoginHistoryService) GetLoginHistory(userID uint, page, pageSize int) ([]domain.LoginHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.historyRepo.GetByUserID(userID, page, pageSize)
}

// CleanupOldRecords deletes login history older than specified days
func (s *LoginHistoryService) CleanupOldRecords(days int) error {
	return s.historyRepo.DeleteOldRecords(days)
}

// parseUserAgent extracts device, browser, and OS from user agent string
func parseUserAgent(userAgent string) (device, browser, os string) {
	ua := strings.ToLower(userAgent)

	// Detect device type
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		if strings.Contains(ua, "android") {
			device = "Android"
		} else if strings.Contains(ua, "iphone") {
			device = "iPhone"
		} else if strings.Contains(ua, "ipad") {
			device = "iPad"
		} else {
			device = "Mobile"
		}
	} else if strings.Contains(ua, "dart") || strings.Contains(ua, "flutter") || strings.Contains(ua, "okhttp") {
		// Dart/Flutter app HTTP clients — cannot determine exact platform from UA alone
		device = "移动端"
	} else {
		device = "Desktop"
	}

	// Detect browser
	if strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/") {
		browser = "Edge"
	} else if strings.Contains(ua, "opr/") || strings.Contains(ua, "opera") {
		browser = "Opera"
	} else if strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox/") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") {
		browser = "Safari"
	} else if strings.Contains(ua, "msie") || strings.Contains(ua, "trident/") {
		browser = "Internet Explorer"
	} else if strings.Contains(ua, "dart") || strings.Contains(ua, "flutter") || strings.Contains(ua, "okhttp") {
		// App client, not a browser
		browser = "App"
	} else {
		browser = "Unknown"
	}

	// Detect OS
	if strings.Contains(ua, "windows nt 10.0") {
		os = "Windows 10"
	} else if strings.Contains(ua, "windows nt 6.3") {
		os = "Windows 8.1"
	} else if strings.Contains(ua, "windows nt 6.2") {
		os = "Windows 8"
	} else if strings.Contains(ua, "windows nt 6.1") {
		os = "Windows 7"
	} else if strings.Contains(ua, "windows") {
		os = "Windows"
	} else if strings.Contains(ua, "mac os x") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") && !strings.Contains(ua, "android") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	} else {
		os = "Unknown"
	}

	return device, browser, os
}