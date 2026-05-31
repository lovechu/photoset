package service

import (
	"fmt"
	"photoset/internal/domain"
	"photoset/internal/repository"
	"time"

	"gorm.io/gorm"
)

// UserLevelService handles user level, achievements, and points mall
type UserLevelService struct {
	levelRepo      *repository.UserLevelRepository
	achievementRepo *repository.AchievementRepository
	pointsMallRepo  *repository.PointsMallRepository
	pointRepo       *repository.UserPointRepository
	db              *gorm.DB
}

// NewUserLevelService creates a new UserLevelService
func NewUserLevelService(
	levelRepo *repository.UserLevelRepository,
	achievementRepo *repository.AchievementRepository,
	pointsMallRepo *repository.PointsMallRepository,
	pointRepo *repository.UserPointRepository,
	db *gorm.DB,
) *UserLevelService {
	return &UserLevelService{
		levelRepo:       levelRepo,
		achievementRepo: achievementRepo,
		pointsMallRepo:  pointsMallRepo,
		pointRepo:       pointRepo,
		db:              db,
	}
}

// Initialize initializes default data
func (s *UserLevelService) Initialize() error {
	if err := s.levelRepo.InitDefaultLevelConfigs(); err != nil {
		return fmt.Errorf("failed to init level configs: %w", err)
	}
	if err := s.achievementRepo.InitDefaultAchievements(); err != nil {
		return fmt.Errorf("failed to init achievements: %w", err)
	}
	if err := s.pointsMallRepo.InitDefaultItems(); err != nil {
		return fmt.Errorf("failed to init points mall items: %w", err)
	}
	return nil
}

// ===== Level System =====

// GetUserLevelInfo returns comprehensive user level information
func (s *UserLevelService) GetUserLevelInfo(userID uint) (map[string]interface{}, error) {
	// Get user points
	userPoint, err := s.pointRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get level config
	levelConfig, err := s.levelRepo.GetLevelConfig(userPoint.Level)
	if err != nil {
		// Use default if not found
		levelConfig = &domain.UserLevelConfig{
			Level:       userPoint.Level,
			Name:        domain.GetLevelName(userPoint.Level),
			MinPoints:   0,
			MaxPoints:   999999,
		}
	}

	// Get next level config
	var nextLevelConfig *domain.UserLevelConfig
	if userPoint.Level < 10 {
		nextLevelConfig, _ = s.levelRepo.GetLevelConfig(userPoint.Level + 1)
	}

	// Calculate progress to next level
	progress := 0.0
	pointsToNext := 0
	if nextLevelConfig != nil {
		currentLevelPoints := levelConfig.MinPoints
		nextLevelPoints := nextLevelConfig.MinPoints
		if nextLevelPoints > currentLevelPoints {
			progress = float64(userPoint.Points-currentLevelPoints) / float64(nextLevelPoints-currentLevelPoints) * 100
			pointsToNext = nextLevelPoints - userPoint.Points
		}
	}

	// Get achievement count
	achievementCount, _ := s.achievementRepo.GetUserAchievementCount(userID)

	return map[string]interface{}{
		"level":            userPoint.Level,
		"level_name":       levelConfig.Name,
		"level_icon":       levelConfig.Icon,
		"level_color":      levelConfig.Color,
		"points":           userPoint.Points,
		"progress":         progress,
		"points_to_next":   pointsToNext,
		"achievement_count": achievementCount,
		"privileges": map[string]interface{}{
			"can_create_post":    levelConfig.CanCreatePost,
			"can_create_reply":   levelConfig.CanCreateReply,
			"can_upload_image":   levelConfig.CanUploadImage,
			"can_upload_video":   levelConfig.CanUploadVideo,
			"can_create_topic":   levelConfig.CanCreateTopic,
			"can_pin_post":       levelConfig.CanPinPost,
			"can_delete_reply":   levelConfig.CanDeleteReply,
			"max_post_per_day":   levelConfig.MaxPostPerDay,
			"max_reply_per_day":  levelConfig.MaxReplyPerDay,
			"max_image_per_post": levelConfig.MaxImagePerPost,
			"max_video_per_post": levelConfig.MaxVideoPerPost,
			"max_post_length":    levelConfig.MaxPostLength,
		},
	}, nil
}

// GetAllLevelConfigs returns all level configurations
func (s *UserLevelService) GetAllLevelConfigs() ([]domain.UserLevelConfig, error) {
	return s.levelRepo.GetAllLevelConfigs()
}

// ===== Achievement System =====

// GetUserAchievements returns user's achievements with unlock status
func (s *UserLevelService) GetUserAchievements(userID uint) ([]map[string]interface{}, error) {
	allAchievements, err := s.achievementRepo.GetAllAchievements()
	if err != nil {
		return nil, err
	}

	unlockedAchievements, err := s.achievementRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}

	// Create map of unlocked achievements
	unlockedMap := make(map[uint]time.Time)
	for _, ua := range unlockedAchievements {
		unlockedMap[ua.AchievementID] = ua.CreatedAt
	}

	// Build response
	var result []map[string]interface{}
	for _, achievement := range allAchievements {
		item := map[string]interface{}{
			"id":             achievement.ID,
			"name":           achievement.Name,
			"title":          achievement.Title,
			"description":    achievement.Description,
			"icon":           achievement.Icon,
			"type":           achievement.Type,
			"condition_type": achievement.ConditionType,
			"condition_value": achievement.ConditionValue,
			"reward_points":  achievement.RewardPoints,
			"is_unlocked":    false,
		}

		if unlockedAt, ok := unlockedMap[achievement.ID]; ok {
			item["is_unlocked"] = true
			item["unlocked_at"] = unlockedAt
		}

		result = append(result, item)
	}

	return result, nil
}

// CheckAndUnlockAchievements checks and unlocks achievements for a user based on their stats
func (s *UserLevelService) CheckAndUnlockAchievements(userID uint, stats map[string]int) ([]domain.Achievement, error) {
	allAchievements, err := s.achievementRepo.GetAllAchievementsIncludingHidden()
	if err != nil {
		return nil, err
	}

	var unlocked []domain.Achievement

	for _, achievement := range allAchievements {
		// Skip if already unlocked
		if s.achievementRepo.HasAchievement(userID, achievement.ID) {
			continue
		}

		// Check condition
		shouldUnlock := false
		switch achievement.ConditionType {
		case "post_count":
			shouldUnlock = stats["post_count"] >= achievement.ConditionValue
		case "reply_count":
			shouldUnlock = stats["reply_count"] >= achievement.ConditionValue
		case "like_received":
			shouldUnlock = stats["like_received"] >= achievement.ConditionValue
		case "following_count":
			shouldUnlock = stats["following_count"] >= achievement.ConditionValue
		case "follower_count":
			shouldUnlock = stats["follower_count"] >= achievement.ConditionValue
		case "level_reached":
			shouldUnlock = stats["level"] >= achievement.ConditionValue
		}

		if shouldUnlock {
			// Unlock achievement
			if err := s.achievementRepo.UnlockAchievement(userID, achievement.ID); err != nil {
				continue
			}

			// Award reward points
			if achievement.RewardPoints > 0 {
				s.pointRepo.AddPoints(userID, achievement.RewardPoints)
				s.pointRepo.LogPointChange(userID, achievement.RewardPoints, "achievement:"+achievement.Name, achievement.ID)
			}

			unlocked = append(unlocked, achievement)
		}
	}

	return unlocked, nil
}

// ===== Points Mall =====

// GetPointsMallItems returns all available items
func (s *UserLevelService) GetPointsMallItems(category string) ([]domain.PointsMallItem, error) {
	if category != "" {
		return s.pointsMallRepo.GetItemsByCategory(category)
	}
	return s.pointsMallRepo.GetAllItems()
}

// ExchangeItem exchanges points for an item
func (s *UserLevelService) ExchangeItem(userID uint, itemID uint) (*domain.UserPointsExchange, error) {
	// Get item
	item, err := s.pointsMallRepo.GetItemByID(itemID)
	if err != nil {
		return nil, domain.ErrItemNotFound
	}

	if !item.IsActive {
		return nil, domain.ErrItemNotActive
	}

	// Get user points
	userPoint, err := s.pointRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Check if can exchange
	if !item.CanExchange(userPoint.Level, userPoint.Points) {
		if item.MinLevel > userPoint.Level {
			return nil, domain.ErrLevelTooLow
		}
		if item.PointsCost > userPoint.Points {
			return nil, domain.ErrInsufficientPoints
		}
		if !item.IsUnlimited && item.GetRemainingStock() <= 0 {
			return nil, domain.ErrOutOfStock
		}
	}

	// Perform exchange
	if err := s.pointsMallRepo.ExchangeItem(userID, item, userPoint); err != nil {
		return nil, err
	}

	// Log the point change
	s.pointRepo.LogPointChange(userID, -item.PointsCost, "exchange:"+item.Name, item.ID)

	// Get the exchange record
	exchanges, _, _ := s.pointsMallRepo.GetUserExchanges(userID, 1, 1)
	if len(exchanges) > 0 {
		return &exchanges[0], nil
	}

	return &domain.UserPointsExchange{
		UserID:    userID,
		ItemID:    itemID,
		Points:    item.PointsCost,
		ItemName:  item.Name,
		ItemType:  item.ItemType,
		ItemValue: item.ItemValue,
		Status:    "completed",
	}, nil
}

// GetUserExchangeHistory returns user's exchange history
func (s *UserLevelService) GetUserExchangeHistory(userID uint, page, pageSize int) ([]domain.UserPointsExchange, int64, error) {
	return s.pointsMallRepo.GetUserExchanges(userID, page, pageSize)
}

// GetPointsMallCategories returns available categories
func (s *UserLevelService) GetPointsMallCategories() []map[string]interface{} {
	return []map[string]interface{}{
		{"key": "badge", "name": "徽章", "icon": "🏅"},
		{"key": "title", "name": "称号", "icon": "👑"},
		{"key": "privilege", "name": "特权", "icon": "⭐"},
		{"key": "virtual_gift", "name": "虚拟礼物", "icon": "🎁"},
	}
}
