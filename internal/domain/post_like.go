package domain

import (
	"time"
)

// PostLike represents a like on a post
type PostLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_post" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_post;index:idx_post_id" json:"post_id"`

	// Associations
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Post      Post      `gorm:"foreignKey:PostID" json:"-"`
}

// TableName specifies the table name
func (PostLike) TableName() string {
	return "post_likes"
}
