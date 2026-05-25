package service

import (
	"errors"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

type FollowService interface {
	Follow(userID, followingID uint) error
	Unfollow(userID, followingID uint) error
	IsFollowing(userID, followingID uint) (bool, error)
	GetFollowingList(userID uint, page, pageSize int) ([]domain.User, int64, error)
	GetFollowerList(userID uint, page, pageSize int) ([]domain.User, int64, error)
	BatchCheckFollowing(userID uint, targetIDs []uint) (map[uint]bool, error)
}

type followService struct {
	followRepo *repository.FollowRepository
	userRepo   repository.UserRepository
}

func NewFollowService(followRepo *repository.FollowRepository, userRepo repository.UserRepository) FollowService {
	return &followService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *followService) Follow(userID, followingID uint) error {
	if userID == followingID {
		return errors.New("不能关注自己")
	}

	// 检查目标用户是否存在
	user, err := s.userRepo.FindByID(followingID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 检查是否已经关注
	exists, err := s.followRepo.Exists(userID, followingID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("已经关注了该用户")
	}

	// 创建关注关系
	follow := &domain.Follow{
		UserID:      userID,
		FollowingID: followingID,
	}
	if err := s.followRepo.Create(follow); err != nil {
		return err
	}

	// 更新计数器
	s.updateCounters(userID, followingID)
	return nil
}

func (s *followService) Unfollow(userID, followingID uint) error {
	exists, err := s.followRepo.Exists(userID, followingID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("未关注该用户")
	}

	if err := s.followRepo.Delete(userID, followingID); err != nil {
		return err
	}

	// 更新计数器
	s.updateCounters(userID, followingID)
	return nil
}

func (s *followService) IsFollowing(userID, followingID uint) (bool, error) {
	return s.followRepo.Exists(userID, followingID)
}

func (s *followService) GetFollowingList(userID uint, page, pageSize int) ([]domain.User, int64, error) {
	users, err := s.followRepo.FindFollowing(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.followRepo.CountFollowing(userID)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *followService) GetFollowerList(userID uint, page, pageSize int) ([]domain.User, int64, error) {
	users, err := s.followRepo.FindFollowers(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.followRepo.CountFollowers(userID)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *followService) BatchCheckFollowing(userID uint, targetIDs []uint) (map[uint]bool, error) {
	return s.followRepo.BatchCheckFollowing(userID, targetIDs)
}

// updateCounters 更新关注/粉丝计数器
func (s *followService) updateCounters(userID, followingID uint) {
	followingCount, _ := s.followRepo.CountFollowing(userID)
	followerCount, _ := s.followRepo.CountFollowers(followingID)

	// 更新关注者的 following_count
	s.userRepo.UpdateField(userID, "following_count", followingCount)
	// 更新被关注者的 follower_count
	s.userRepo.UpdateField(followingID, "follower_count", followerCount)
}
