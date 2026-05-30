package domain

import (
	"time"
)

// PostFavorite represents a user's favorite post
type PostFavorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_post" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_post" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName specifies the table name
func (PostFavorite) TableName() string {
	return "post_favorites"
}