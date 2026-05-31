package domain

import (
	"time"

	"gorm.io/gorm"
)

// UserLevelConfig defines the configuration for each user level
type UserLevelConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Level       int    `gorm:"uniqueIndex;not null" json:"level"`
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	Icon        string `gorm:"type:varchar(255);default:''" json:"icon"`
	Color       string `gorm:"type:varchar(20);default:'#FFD700'" json:"color"`
	MinPoints   int    `gorm:"not null;default:0" json:"min_points"`
	MaxPoints   int    `gorm:"not null;default:0" json:"max_points"`
	Description string `gorm:"type:text" json:"description"`
	
	// Privileges unlocked at this level
	CanCreatePost    bool `gorm:"default:true" json:"can_create_post"`
	CanCreateReply   bool `gorm:"default:true" json:"can_create_reply"`
	CanUploadImage   bool `gorm:"default:true" json:"can_upload_image"`
	CanUploadVideo   bool `gorm:"default:false" json:"can_upload_video"`
	CanCreateTopic   bool `gorm:"default:false" json:"can_create_topic"`
	CanPinPost       bool `gorm:"default:false" json:"can_pin_post"`
	CanDeleteReply   bool `gorm:"default:false" json:"can_delete_reply"`
	MaxPostPerDay    int  `gorm:"default:5" json:"max_post_per_day"`
	MaxReplyPerDay   int  `gorm:"default:10" json:"max_reply_per_day"`
	MaxImagePerPost  int  `gorm:"default:9" json:"max_image_per_post"`
	MaxVideoPerPost  int  `gorm:"default:0" json:"max_video_per_post"`
	MaxPostLength    int  `gorm:"default:5000" json:"max_post_length"`
	
	// Reward for reaching this level
	RewardPoints    int    `gorm:"default:0" json:"reward_points"`
	RewardBadge     string `gorm:"type:varchar(50);default:''" json:"reward_badge"`
	RewardTitle     string `gorm:"type:varchar(50);default:''" json:"reward_title"`
}

// TableName specifies the table name
func (UserLevelConfig) TableName() string {
	return "user_level_configs"
}

// DefaultLevelConfigs returns default level configurations
func DefaultLevelConfigs() []UserLevelConfig {
	return []UserLevelConfig{
		{
			Level:       1,
			Name:        "新手",
			Icon:        "🌱",
			Color:       "#8BC34A",
			MinPoints:   0,
			MaxPoints:   99,
			Description: "刚加入社区的新手",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   false,
			CanCreateTopic:   false,
			CanPinPost:       false,
			CanDeleteReply:   false,
			MaxPostPerDay:    5,
			MaxReplyPerDay:   10,
			MaxImagePerPost:  9,
			MaxVideoPerPost:  0,
			MaxPostLength:    5000,
			RewardPoints:     0,
			RewardBadge:      "newbie",
			RewardTitle:      "社区新人",
		},
		{
			Level:       2,
			Name:        "活跃成员",
			Icon:        "⭐",
			Color:       "#2196F3",
			MinPoints:   100,
			MaxPoints:   499,
			Description: "积极参与社区讨论的成员",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   false,
			CanCreateTopic:   false,
			CanPinPost:       false,
			CanDeleteReply:   false,
			MaxPostPerDay:    10,
			MaxReplyPerDay:   20,
			MaxImagePerPost:  9,
			MaxVideoPerPost:  0,
			MaxPostLength:    8000,
			RewardPoints:     50,
			RewardBadge:      "active",
			RewardTitle:      "活跃达人",
		},
		{
			Level:       3,
			Name:        "资深成员",
			Icon:        "🌟",
			Color:       "#9C27B0",
			MinPoints:   500,
			MaxPoints:   1999,
			Description: "社区的资深贡献者",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   false,
			CanCreateTopic:   true,
			CanPinPost:       false,
			CanDeleteReply:   false,
			MaxPostPerDay:    15,
			MaxReplyPerDay:   30,
			MaxImagePerPost:  12,
			MaxVideoPerPost:  0,
			MaxPostLength:    10000,
			RewardPoints:     100,
			RewardBadge:      "senior",
			RewardTitle:      "资深社区成员",
		},
		{
			Level:       4,
			Name:        "金牌成员",
			Icon:        "🏅",
			Color:       "#FF9800",
			MinPoints:   2000,
			MaxPoints:   4999,
			Description: "社区的金牌贡献者",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       false,
			CanDeleteReply:   false,
			MaxPostPerDay:    20,
			MaxReplyPerDay:   50,
			MaxImagePerPost:  15,
			MaxVideoPerPost:  1,
			MaxPostLength:    15000,
			RewardPoints:     200,
			RewardBadge:      "gold",
			RewardTitle:      "金牌创作者",
		},
		{
			Level:       5,
			Name:        "钻石成员",
			Icon:        "💎",
			Color:       "#00BCD4",
			MinPoints:   5000,
			MaxPoints:   9999,
			Description: "社区的钻石级成员",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    30,
			MaxReplyPerDay:   100,
			MaxImagePerPost:  20,
			MaxVideoPerPost:  3,
			MaxPostLength:    20000,
			RewardPoints:     500,
			RewardBadge:      "diamond",
			RewardTitle:      "钻石创作者",
		},
		{
			Level:       6,
			Name:        "至尊成员",
			Icon:        "👑",
			Color:       "#F44336",
			MinPoints:   10000,
			MaxPoints:   19999,
			Description: "社区的至尊级成员",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    50,
			MaxReplyPerDay:   200,
			MaxImagePerPost:  30,
			MaxVideoPerPost:  5,
			MaxPostLength:    30000,
			RewardPoints:     1000,
			RewardBadge:      "supreme",
			RewardTitle:      "至尊创作者",
		},
		{
			Level:       7,
			Name:        "荣耀L7",
			Icon:        "🏆",
			Color:       "#E91E63",
			MinPoints:   20000,
			MaxPoints:   29999,
			Description: "荣耀等级第七级",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    100,
			MaxReplyPerDay:   500,
			MaxImagePerPost:  50,
			MaxVideoPerPost:  10,
			MaxPostLength:    50000,
			RewardPoints:     2000,
			RewardBadge:      "glory7",
			RewardTitle:      "荣耀守护者",
		},
		{
			Level:       8,
			Name:        "荣耀L8",
			Icon:        "🏆",
			Color:       "#673AB7",
			MinPoints:   30000,
			MaxPoints:   39999,
			Description: "荣耀等级第八级",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    100,
			MaxReplyPerDay:   500,
			MaxImagePerPost:  50,
			MaxVideoPerPost:  10,
			MaxPostLength:    50000,
			RewardPoints:     3000,
			RewardBadge:      "glory8",
			RewardTitle:      "荣耀精英",
		},
		{
			Level:       9,
			Name:        "荣耀L9",
			Icon:        "🏆",
			Color:       "#3F51B5",
			MinPoints:   40000,
			MaxPoints:   49999,
			Description: "荣耀等级第九级",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    100,
			MaxReplyPerDay:   500,
			MaxImagePerPost:  50,
			MaxVideoPerPost:  10,
			MaxPostLength:    50000,
			RewardPoints:     5000,
			RewardBadge:      "glory9",
			RewardTitle:      "荣耀大师",
		},
		{
			Level:       10,
			Name:        "荣耀L10",
			Icon:        "🏆",
			Color:       "#FFD700",
			MinPoints:   50000,
			MaxPoints:   999999,
			Description: "荣耀等级最高级",
			CanCreatePost:    true,
			CanCreateReply:   true,
			CanUploadImage:   true,
			CanUploadVideo:   true,
			CanCreateTopic:   true,
			CanPinPost:       true,
			CanDeleteReply:   true,
			MaxPostPerDay:    999,
			MaxReplyPerDay:   999,
			MaxImagePerPost:  99,
			MaxVideoPerPost:  20,
			MaxPostLength:    99999,
			RewardPoints:     10000,
			RewardBadge:      "glory10",
			RewardTitle:      "荣耀传说",
		},
	}
}
