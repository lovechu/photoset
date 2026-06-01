package domain

import (
	"time"
)

// PostShare represents a share record for a post
type PostShare struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uint   `gorm:"not null;index:idx_share_user" json:"user_id"`
	PostID    uint   `gorm:"not null;index:idx_share_post" json:"post_id"`
	Platform  string `gorm:"type:varchar(20);not null;default:'other'" json:"platform"` // wechat, weibo, link, other

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
	Post Post `gorm:"foreignKey:PostID" json:"-"`
}

// TableName specifies the table name
func (PostShare) TableName() string {
	return "post_shares"
}
