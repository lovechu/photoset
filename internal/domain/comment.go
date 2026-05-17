package domain

import "time"

// Comment 评论模型
type Comment struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	PhotoSetID uint       `gorm:"not null;index" json:"photoset_id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	Content    string     `gorm:"type:text;not null" json:"content"`
	ImageURL   string     `gorm:"size:500;default:''" json:"image_url"`
	ParentID   *uint      `gorm:"index" json:"parent_id"` // 回复的评论ID，nil表示顶级评论
	LikeCount  int        `gorm:"default:0" json:"like_count"`

	// 关联
	User   User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// CommentLike 评论点赞模型
type CommentLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_user_comment" json:"user_id"`
	CommentID uint      `gorm:"not null;uniqueIndex:uk_user_comment" json:"comment_id"`
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}
