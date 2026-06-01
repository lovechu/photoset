package domain

import (
	"time"

	"gorm.io/gorm"
)

// AchievementType defines the type of achievement
type AchievementType string

const (
	AchievementTypePost     AchievementType = "post"     // 发帖相关
	AchievementTypeReply    AchievementType = "reply"    // 回复相关
	AchievementTypeLike     AchievementType = "like"     // 点赞相关
	AchievementTypeFollow   AchievementType = "follow"   // 关注相关
	AchievementTypeLevel    AchievementType = "level"    // 等级相关
	AchievementTypeSpecial  AchievementType = "special"  // 特殊成就
)

// Achievement defines an achievement that users can unlock
type Achievement struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Title       string         `gorm:"type:varchar(100);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(255);default:''" json:"icon"`
	BadgeImage  string         `gorm:"type:varchar(255);default:''" json:"badge_image"`
	Type        AchievementType `gorm:"type:varchar(20);not null;index" json:"type"`
	
	// Unlock condition
	ConditionType  string `gorm:"type:varchar(50);not null" json:"condition_type"`  // e.g., "post_count", "reply_count", "like_received", "level_reached"
	ConditionValue int    `gorm:"not null;default:0" json:"condition_value"`
	
	// Reward
	RewardPoints int    `gorm:"default:0" json:"reward_points"`
	RewardTitle  string `gorm:"type:varchar(50);default:''" json:"reward_title"`
	
	// Sort order and visibility
	SortOrder int  `gorm:"default:0" json:"sort_order"`
	IsHidden  bool `gorm:"default:false" json:"is_hidden"`
}

// TableName specifies the table name
func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement represents a user's unlocked achievement
type UserAchievement struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	AchievementID uint      `gorm:"not null;index" json:"achievement_id"`
	
	// Associations
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Achievement Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
}

// TableName specifies the table name
func (UserAchievement) TableName() string {
	return "user_achievements"
}

// DefaultAchievements returns default achievements
func DefaultAchievements() []Achievement {
	return []Achievement{
		// Post achievements
		{
			Name:           "first_post",
			Title:          "初次发帖",
			Description:    "发布第一篇帖子",
			Icon:           "📝",
			Type:           AchievementTypePost,
			ConditionType:  "post_count",
			ConditionValue: 1,
			RewardPoints:   20,
			RewardTitle:    "初出茅庐",
			SortOrder:      1,
		},
		{
			Name:           "post_10",
			Title:          "小有成就",
			Description:    "累计发布10篇帖子",
			Icon:           "📚",
			Type:           AchievementTypePost,
			ConditionType:  "post_count",
			ConditionValue: 10,
			RewardPoints:   50,
			SortOrder:      2,
		},
		{
			Name:           "post_50",
			Title:          "创作达人",
			Description:    "累计发布50篇帖子",
			Icon:           "✍️",
			Type:           AchievementTypePost,
			ConditionType:  "post_count",
			ConditionValue: 50,
			RewardPoints:   200,
			SortOrder:      3,
		},
		{
			Name:           "post_100",
			Title:          "百帖大神",
			Description:    "累计发布100篇帖子",
			Icon:           "🎯",
			Type:           AchievementTypePost,
			ConditionType:  "post_count",
			ConditionValue: 100,
			RewardPoints:   500,
			SortOrder:      4,
		},
		
		// Reply achievements
		{
			Name:           "first_reply",
			Title:          "初次回复",
			Description:    "发布第一篇回复",
			Icon:           "💬",
			Type:           AchievementTypeReply,
			ConditionType:  "reply_count",
			ConditionValue: 1,
			RewardPoints:   10,
			RewardTitle:    "热心网友",
			SortOrder:      10,
		},
		{
			Name:           "reply_100",
			Title:          "评论达人",
			Description:    "累计发布100篇回复",
			Icon:           "🗣️",
			Type:           AchievementTypeReply,
			ConditionType:  "reply_count",
			ConditionValue: 100,
			RewardPoints:   100,
			SortOrder:      11,
		},
		
		// Like achievements
		{
			Name:           "first_like",
			Title:          "初次获赞",
			Description:    "收到第一个点赞",
			Icon:           "❤️",
			Type:           AchievementTypeLike,
			ConditionType:  "like_received",
			ConditionValue: 1,
			RewardPoints:   10,
			SortOrder:      20,
		},
		{
			Name:           "like_100",
			Title:          "万人迷",
			Description:    "累计收到100个点赞",
			Icon:           "💖",
			Type:           AchievementTypeLike,
			ConditionType:  "like_received",
			ConditionValue: 100,
			RewardPoints:   100,
			SortOrder:      21,
		},
		{
			Name:           "like_1000",
			Title:          "超级网红",
			Description:    "累计收到1000个点赞",
			Icon:           "🌟",
			Type:           AchievementTypeLike,
			ConditionType:  "like_received",
			ConditionValue: 1000,
			RewardPoints:   500,
			SortOrder:      22,
		},
		
		// Follow achievements
		{
			Name:           "first_follow",
			Title:          "初次关注",
			Description:    "关注第一个用户",
			Icon:           "🤝",
			Type:           AchievementTypeFollow,
			ConditionType:  "following_count",
			ConditionValue: 1,
			RewardPoints:   10,
			SortOrder:      30,
		},
		{
			Name:           "follower_10",
			Title:          "小有人气",
			Description:    "拥有10个粉丝",
			Icon:           "👥",
			Type:           AchievementTypeFollow,
			ConditionType:  "follower_count",
			ConditionValue: 10,
			RewardPoints:   50,
			SortOrder:      31,
		},
		{
			Name:           "follower_100",
			Title:          "人气之星",
			Description:    "拥有100个粉丝",
			Icon:           "⭐",
			Type:           AchievementTypeFollow,
			ConditionType:  "follower_count",
			ConditionValue: 100,
			RewardPoints:   200,
			SortOrder:      32,
		},
		
		// Level achievements
		{
			Name:           "level_5",
			Title:          "钻石之路",
			Description:    "达到5级",
			Icon:           "💎",
			Type:           AchievementTypeLevel,
			ConditionType:  "level_reached",
			ConditionValue: 5,
			RewardPoints:   300,
			SortOrder:      40,
		},
		{
			Name:           "level_10",
			Title:          "荣耀巅峰",
			Description:    "达到10级",
			Icon:           "🏆",
			Type:           AchievementTypeLevel,
			ConditionType:  "level_reached",
			ConditionValue: 10,
			RewardPoints:   1000,
			SortOrder:      41,
		},
		
		// Special achievements
		{
			Name:           "early_bird",
			Title:          "早起的鸟儿",
			Description:    "在早上6点前发帖",
			Icon:           "🐦",
			Type:           AchievementTypeSpecial,
			ConditionType:  "special",
			ConditionValue: 1,
			RewardPoints:   50,
			IsHidden:       true,
			SortOrder:      50,
		},
		{
			Name:           "night_owl",
			Title:          "夜猫子",
			Description:    "在凌晨2点后发帖",
			Icon:           "🦉",
			Type:           AchievementTypeSpecial,
			ConditionType:  "special",
			ConditionValue: 2,
			RewardPoints:   50,
			IsHidden:       true,
			SortOrder:      51,
		},
	}
}
