package domain

import (
	"time"
)

// UserPrivacySetting represents user privacy settings
type UserPrivacySetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint `gorm:"not null;uniqueIndex" json:"user_id"`

	// 个人资料可见性
	ShowProfile bool `gorm:"default:true" json:"show_profile"` // 是否展示个人资料
	ShowPosts   bool `gorm:"default:true" json:"show_posts"`   // 是否展示帖子
	ShowFavorites bool `gorm:"default:true" json:"show_favorites"` // 是否展示收藏

	// 搜索相关
	AllowSearch bool `gorm:"default:true" json:"allow_search"` // 是否允许被搜索

	// 消息相关
	AllowMessage bool `gorm:"default:true" json:"allow_message"` // 是否允许接收私信

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (UserPrivacySetting) TableName() string {
	return "user_privacy_settings"
}
