package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
)

// PointService handles user points and level calculations
type PointService struct {
	pointRepo *repository.UserPointRepository
}

// NewPointService creates a new PointService
func NewPointService(pointRepo *repository.UserPointRepository) *PointService {
	return &PointService{pointRepo: pointRepo}
}

// GetUserPoints gets user points and level
func (s *PointService) GetUserPoints(userID uint) (*domain.UserPoint, error) {
	return s.pointRepo.GetByUserID(userID)
}

// AddPointsForPost adds points for creating a post (with daily limit check)
func (s *PointService) AddPointsForPost(userID uint) error {
	// Check daily limit (max 50 points = 5 posts)
	todayPoints, err := s.pointRepo.GetTodayPostPoints(userID)
	if err != nil {
		return err
	}

	if todayPoints >= 50 {
		return domain.ErrDailyLimitReached
	}

	// Add points
	if err := s.pointRepo.AddPoints(userID, 10); err != nil {
		return err
	}

	// Log the change
	return s.pointRepo.LogPointChange(userID, 10, "post_create", 0)
}

// AddPointsForReply adds points for creating a reply (with daily limit check)
func (s *PointService) AddPointsForReply(userID uint) error {
	// Check daily limit (max 30 points = 6 replies)
	todayPoints, err := s.pointRepo.GetTodayReplyPoints(userID)
	if err != nil {
		return err
	}

	if todayPoints >= 30 {
		return domain.ErrDailyLimitReached
	}

	// Add points
	if err := s.pointRepo.AddPoints(userID, 5); err != nil {
		return err
	}

	// Log the change
	return s.pointRepo.LogPointChange(userID, 5, "reply_create", 0)
}

// AddPointsForLiked adds points when user's post/reply is liked
func (s *PointService) AddPointsForLiked(userID uint, points int) error {
	if err := s.pointRepo.AddPoints(userID, points); err != nil {
		return err
	}

	// Log the change
	action := "post_liked"
	if points == 1 {
		action = "reply_liked"
	}
	return s.pointRepo.LogPointChange(userID, points, action, 0)
}

// DeductPointsForDelete deducts points when post/reply is deleted by admin
func (s *PointService) DeductPointsForDelete(userID uint) error {
	if err := s.pointRepo.AddPoints(userID, -20); err != nil {
		return err
	}

	// Log the change
	return s.pointRepo.LogPointChange(userID, -20, "admin_delete", 0)
}

// AdjustPoints manually adjusts user points (admin function)
func (s *PointService) AdjustPoints(userID uint, points int, reason string) error {
	if err := s.pointRepo.AddPoints(userID, points); err != nil {
		return err
	}

	// Log the change
	return s.pointRepo.LogPointChange(userID, points, "admin_adjust:"+reason, 0)
}

// GetLevelInfo gets user level information
func (s *PointService) GetLevelInfo(userID uint) (level int, levelName string, currentPoints int, nextLevelPoints int, err error) {
	userPoint, err := s.pointRepo.GetByUserID(userID)
	if err != nil {
		return 0, "", 0, 0, err
	}

	level = userPoint.Level
	levelName = domain.GetLevelName(level)
	currentPoints = userPoint.Points
	nextLevelPoints = userPoint.GetNextLevelPoints()

	return level, levelName, currentPoints, nextLevelPoints, nil
}


