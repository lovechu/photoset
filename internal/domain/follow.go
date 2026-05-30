package domain

import (
	"time"
)

// Follow represents a user following another user
type Follow struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID     uint `gorm:"not null;uniqueIndex:idx_user_follow" json:"user_id"`     // 关注者
	FollowingID uint `gorm:"not null;uniqueIndex:idx_user_follow;index:idx_following_id" json:"following_id"` // 被关注者

	// Associations
	User      User `gorm:"foreignKey:UserID" json:"-"`
	Following User `gorm:"foreignKey:FollowingID" json:"-"`
}

// TableName specifies the table name
func (Follow) TableName() string {
	return "follows"
}
