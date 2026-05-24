package domain

import (
	"time"
)

// PostReply represents a reply to a post (supports nested replies)
type PostReply struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	PostID         uint       `gorm:"not null;index" json:"post_id"`
	UserID         uint       `gorm:"not null;index" json:"user_id"`
	Content        string     `gorm:"type:text;not null" json:"content"`
	ParentReplyID  *uint      `gorm:"index" json:"parent_reply_id"` // nil means direct reply to post
	LikeCount      int        `gorm:"not null;default:0" json:"like_count"`

	// Associations
	User           User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post           Post          `gorm:"foreignKey:PostID" json:"-"`
	ParentReply    *PostReply    `gorm:"foreignKey:ParentReplyID" json:"-"`
	Children       []PostReply   `gorm:"foreignKey:ParentReplyID" json:"children,omitempty"`
	Likes          []PostReplyLike `gorm:"foreignKey:ReplyID" json:"-"`
}

// TableName specifies the table name
func (PostReply) TableName() string {
	return "post_replies"
}

// Validate checks if the reply data is valid
func (r *PostReply) Validate() error {
	if r.Content == "" {
		return ErrReplyContentRequired
	}
	if len(r.Content) > 2000 {
		return ErrReplyContentTooLong
	}
	return nil
}

// IsDirectReply returns true if this is a direct reply to a post (not a reply to another reply)
func (r *PostReply) IsDirectReply() bool {
	return r.ParentReplyID == nil
}

// GetChildren returns the nested replies (children)
func (r *PostReply) GetChildren() []PostReply {
	return r.Children
}
