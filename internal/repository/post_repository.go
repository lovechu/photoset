package repository

import (
	"photoset/internal/domain"
	"time"

	"gorm.io/gorm"
)

// PostRepository handles database operations for posts
type PostRepository struct {
	DB *gorm.DB
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// Create creates a new post
func (r *PostRepository) Create(post *domain.Post) error {
	return r.DB.Create(post).Error
}

// FindByID finds a post by ID with associations
func (r *PostRepository) FindByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.DB.Preload("User").Preload("Photoset").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindByIDWithCounts finds a post by ID and loads counts
func (r *PostRepository) FindByIDWithCounts(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.DB.Preload("User").Preload("Photoset").First(&post, id).Error
	if err != nil {
		return nil, err
	}

	// Load replies count and likes count
	var replyCount, likeCount int64
	r.DB.Model(&domain.PostReply{}).Where("post_id = ?", id).Count(&replyCount)
	r.DB.Model(&domain.PostLike{}).Where("post_id = ?", id).Count(&likeCount)

	post.ReplyCount = int(replyCount)
	post.LikeCount = int(likeCount)

	return &post, nil
}

// List returns a list of posts with filtering and pagination
func (r *PostRepository) List(page, pageSize int, category, visibility string, userRole string, userID uint, orderBy string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	query := r.DB.Model(&domain.Post{}).Where("status = ?", "approved")

	// Visibility filter based on user role
	if userID == 0 {
		// Not logged in, only public
		query = query.Where("visibility = ?", "public")
	} else if userRole == "admin" {
		// Admin can see all
	} else if userRole == "vip" || userRole == "creator" {
		// VIP/Creator can see public, member, vip
		query = query.Where("visibility IN (?, ?, ?)", "public", "member", "vip")
	} else if userRole == "member" {
		// Member can see public, member
		query = query.Where("visibility IN (?, ?)", "public", "member")
	} else {
		// Regular user can only see public
		query = query.Where("visibility = ?", "public")
	}

	// Category filter
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ordering
	switch orderBy {
	case "hot":
		// Hot posts: (reply_count * 2 + like_count + view_count / 10)
		// Only for posts within 7 days
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		query = query.Where("created_at >= ?", sevenDaysAgo)
		query = query.Order("reply_count * 2 + like_count + view_count / 10 DESC, created_at DESC")
	case "reply_count":
		query = query.Order("reply_count DESC, created_at DESC")
	case "like_count":
		query = query.Order("like_count DESC, created_at DESC")
	default:
		// Default: pinned first, then latest
		query = query.Order("is_pinned DESC, created_at DESC")
	}

	// Pagination
	offset := (page - 1) * pageSize
	err := query.Preload("User").Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// ListForAdmin returns all posts for admin (including non-approved)
func (r *PostRepository) ListForAdmin(page, pageSize int, status string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	query := r.DB.Model(&domain.Post{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Update updates a post
func (r *PostRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateStatus updates a post's status
func (r *PostRepository) UpdateStatus(id uint, status string) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("status", status).Error
}

// TogglePin toggles the pin status of a post
func (r *PostRepository) TogglePin(id uint) error {
	var post domain.Post
	if err := r.DB.First(&post, id).Error; err != nil {
		return err
	}

	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("is_pinned", !post.IsPinned).Error
}

// Delete deletes a post (hard delete)
func (r *PostRepository) Delete(id uint) error {
	return r.DB.Unscoped().Delete(&domain.Post{}, id).Error
}

// IncrementViewCount increments the view count
func (r *PostRepository) IncrementViewCount(id uint) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// IncrementReplyCount increments the reply count
func (r *PostRepository) IncrementReplyCount(id uint) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("reply_count", gorm.Expr("reply_count + 1")).Error
}

// DecrementReplyCount decrements the reply count
func (r *PostRepository) DecrementReplyCount(id uint) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("reply_count", gorm.Expr("GREATEST(reply_count - 1, 0)")).Error
}

// IncrementLikeCount increments the like count
func (r *PostRepository) IncrementLikeCount(id uint) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}

// DecrementLikeCount decrements the like count
func (r *PostRepository) DecrementLikeCount(id uint) error {
	return r.DB.Model(&domain.Post{}).Where("id = ?", id).Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// FindByUserID finds posts by user ID
func (r *PostRepository) FindByUserID(userID uint, page, pageSize int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	query := r.DB.Model(&domain.Post{}).Where("user_id = ? AND status = ?", userID, "approved")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
