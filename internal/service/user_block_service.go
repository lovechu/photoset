package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// UserBlockService handles user block business logic
type UserBlockService struct {
	blockRepo *repository.UserBlockRepository
}

// NewUserBlockService creates a new UserBlockService
func NewUserBlockService(blockRepo *repository.UserBlockRepository) *UserBlockService {
	return &UserBlockService{blockRepo: blockRepo}
}

// BlockUser blocks a user (complete block)
func (s *UserBlockService) BlockUser(userID, blockedID uint) error {
	if userID == blockedID {
		return errors.New("不能拉黑自己")
	}

	// Check if already blocked
	existing, err := s.blockRepo.FindByUserAndBlocked(userID, blockedID)
	if err != nil {
		return err
	}
	if existing != nil {
		// Update block type if different
		if existing.BlockType != "block" {
			existing.BlockType = "block"
			return s.blockRepo.Create(existing) // This will update due to unique constraint
		}
		return errors.New("已经拉黑该用户")
	}

	block := &domain.UserBlock{
		UserID:    userID,
		BlockedID: blockedID,
		BlockType: "block",
	}
	return s.blockRepo.Create(block)
}

// MuteUser mutes a user (hide content but no notification)
func (s *UserBlockService) MuteUser(userID, mutedID uint) error {
	if userID == mutedID {
		return errors.New("不能静音自己")
	}

	// Check if already blocked/muted
	existing, err := s.blockRepo.FindByUserAndBlocked(userID, mutedID)
	if err != nil {
		return err
	}
	if existing != nil {
		// Update to mute if currently blocked
		if existing.BlockType != "mute" {
			existing.BlockType = "mute"
			return s.blockRepo.Create(existing)
		}
		return errors.New("已经静音该用户")
	}

	block := &domain.UserBlock{
		UserID:    userID,
		BlockedID: mutedID,
		BlockType: "mute",
	}
	return s.blockRepo.Create(block)
}

// UnblockUser unblocks or unmutes a user
func (s *UserBlockService) UnblockUser(userID, blockedID uint) error {
	return s.blockRepo.Delete(userID, blockedID)
}

// GetBlockedUsers gets list of blocked users
func (s *UserBlockService) GetBlockedUsers(userID uint, page, pageSize int) ([]domain.UserBlock, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.blockRepo.GetBlockedUsers(userID, page, pageSize)
}

// IsBlocked checks if target user is blocked by userID
func (s *UserBlockService) IsBlocked(userID, targetID uint) (bool, error) {
	return s.blockRepo.IsBlocked(userID, targetID)
}

// IsMuted checks if target user is muted by userID
func (s *UserBlockService) IsMuted(userID, targetID uint) (bool, error) {
	return s.blockRepo.IsMuted(userID, targetID)
}

// GetBlockedUserIDs gets list of blocked user IDs for content filtering
func (s *UserBlockService) GetBlockedUserIDs(userID uint) ([]uint, error) {
	return s.blockRepo.GetBlockedUserIDs(userID)
}

// GetBlockStatus gets the block status between two users
func (s *UserBlockService) GetBlockStatus(userID, targetID uint) (string, error) {
	block, err := s.blockRepo.FindByUserAndBlocked(userID, targetID)
	if err != nil {
		return "", err
	}
	if block == nil {
		return "none", nil
	}
	return block.BlockType, nil
}
