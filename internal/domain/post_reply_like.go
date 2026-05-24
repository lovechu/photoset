package domain

import (
	"time"
)

// PostReplyLike represents a like on a post reply
type PostReplyLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_reply" json:"user_id"`
	ReplyID    uint      `gorm:"not null;uniqueIndex:idx_user_reply;index:idx_reply_id" json:"reply_id"`

	// Associations
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Reply     PostReply `gorm:"foreignKey:ReplyID" json:"-"`
}

// TableName specifies the table name
func (PostReplyLike) TableName() string {
	return "post_reply_likes"
}
