package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// UserPrivacyService handles user privacy settings business logic
type UserPrivacyService struct {
	privacyRepo *repository.UserPrivacyRepository
}

// NewUserPrivacyService creates a new UserPrivacyService
func NewUserPrivacyService(privacyRepo *repository.UserPrivacyRepository) *UserPrivacyService {
	return &UserPrivacyService{privacyRepo: privacyRepo}
}

// GetPrivacySettings gets privacy settings for a user
func (s *UserPrivacyService) GetPrivacySettings(userID uint) (*domain.UserPrivacySetting, error) {
	return s.privacyRepo.GetByUserID(userID)
}

// UpdatePrivacySettings updates privacy settings for a user
func (s *UserPrivacyService) UpdatePrivacySettings(userID uint, setting *domain.UserPrivacySetting) error {
	// 确保设置属于当前用户
	setting.UserID = userID
	
	// 验证并设置默认值
	// 这些布尔字段在Go中默认为false，所以我们需要确保它们有合理的默认值
	// 但这里我们信任前端传入的值，因为它们是显式的
	
	return s.privacyRepo.CreateOrUpdate(setting)
}

// CheckProfileVisibility checks if a user's profile is visible to another user
func (s *UserPrivacyService) CheckProfileVisibility(userID, viewerID uint) (bool, error) {
	// 如果是自己，总是可见
	if userID == viewerID {
		return true, nil
	}
	
	setting, err := s.privacyRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	
	return setting.ShowProfile, nil
}

// CheckPostsVisibility checks if a user's posts are visible to another user
func (s *UserPrivacyService) CheckPostsVisibility(userID, viewerID uint) (bool, error) {
	// 如果是自己，总是可见
	if userID == viewerID {
		return true, nil
	}
	
	setting, err := s.privacyRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	
	return setting.ShowPosts, nil
}

// CheckFavoritesVisibility checks if a user's favorites are visible to another user
func (s *UserPrivacyService) CheckFavoritesVisibility(userID, viewerID uint) (bool, error) {
	// 如果是自己，总是可见
	if userID == viewerID {
		return true, nil
	}
	
	setting, err := s.privacyRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	
	return setting.ShowFavorites, nil
}

// CheckSearchable checks if a user is searchable
func (s *UserPrivacyService) CheckSearchable(userID uint) (bool, error) {
	setting, err := s.privacyRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	
	return setting.AllowSearch, nil
}

// CheckMessageable checks if a user can receive messages
func (s *UserPrivacyService) CheckMessageable(userID uint) (bool, error) {
	setting, err := s.privacyRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}
	
	return setting.AllowMessage, nil
}