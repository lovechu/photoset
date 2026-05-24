package domain

import (
	"time"
)

// PostCategory represents the category of a post
type PostCategory string

const (
	CategoryDiscussion PostCategory = "discussion"
	CategoryQA        PostCategory = "qa"
	CategoryShowcase  PostCategory = "showcase"
	CategorySuggestion PostCategory = "suggestion"
)

// PostVisibility represents the visibility level of a post
type PostVisibility string

const (
	VisibilityPublic  PostVisibility = "public"
	VisibilityMember PostVisibility = "member"
	VisibilityVIP    PostVisibility = "vip"
	VisibilityAdmin  PostVisibility = "admin"
)

// PostStatus represents the status of a post
type PostStatus string

const (
	PostStatusApproved PostStatus = "approved"
	PostStatusPending  PostStatus = "pending"
	PostStatusRejected PostStatus = "rejected"
)

// Post represents a community post
type Post struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `gorm:"index" json:"deleted_at,omitempty"`

	UserID      uint          `gorm:"not null;index" json:"user_id"`
	Title       string        `gorm:"type:varchar(200);not null" json:"title"`
	Content     string        `gorm:"type:text;not null" json:"content"`
	PhotosetID  *uint         `json:"photoset_id"` // nullable, optional association with photoset
	Category    string        `gorm:"type:varchar(20);not null;default:'discussion'" json:"category"`
	Visibility  string        `gorm:"type:varchar(20);not null;default:'public'" json:"visibility"`
	IsPinned   bool          `gorm:"not null;default:false" json:"is_pinned"`
	ViewCount   int           `gorm:"not null;default:0" json:"view_count"`
	ReplyCount  int           `gorm:"not null;default:0" json:"reply_count"`
	LikeCount   int           `gorm:"not null;default:0" json:"like_count"`
	Status      string        `gorm:"type:varchar(20);not null;default:'approved'" json:"status"`

	// Associations
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Photoset  *PhotoSet  `gorm:"foreignKey:PhotosetID" json:"photoset,omitempty"`
	Replies   []PostReply `gorm:"foreignKey:PostID" json:"replies,omitempty"`
	Likes     []PostLike  `gorm:"foreignKey:PostID" json:"-"`
}

// TableName specifies the table name
func (Post) TableName() string {
	return "posts"
}

// Validate checks if the post data is valid
func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrTitleRequired
	}
	if len(p.Title) > 200 {
		return ErrTitleTooLong
	}
	if p.Content == "" {
		return ErrContentRequired
	}
	if len(p.Content) > 5000 {
		return ErrContentTooLong
	}
	// Validate category
	validCategories := []PostCategory{CategoryDiscussion, CategoryQA, CategoryShowcase, CategorySuggestion}
	valid := false
	for _, c := range validCategories {
		if PostCategory(p.Category) == c {
			valid = true
			break
		}
	}
	if !valid && p.Category != "" {
		return ErrInvalidCategory
	}
	// Validate visibility
	validVisibilities := []PostVisibility{VisibilityPublic, VisibilityMember, VisibilityVIP, VisibilityAdmin}
	valid = false
	for _, v := range validVisibilities {
		if PostVisibility(p.Visibility) == v {
			valid = true
			break
		}
	}
	if !valid && p.Visibility != "" {
		return ErrInvalidVisibility
	}
	return nil
}

// IncrementViewCount increments the view count
func (p *Post) IncrementViewCount() {
	p.ViewCount++
}

// UpdateCounts updates reply_count and like_count from related records
func (p *Post) UpdateCounts() {
	p.ReplyCount = len(p.Replies)
	p.LikeCount = len(p.Likes)
}
