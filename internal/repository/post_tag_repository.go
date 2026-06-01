package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostTagRepository handles database operations for post-tag relationships
type PostTagRepository struct {
	DB *gorm.DB
}

// NewPostTagRepository creates a new PostTagRepository
func NewPostTagRepository(db *gorm.DB) *PostTagRepository {
	return &PostTagRepository{DB: db}
}

// AddTagsToPost adds tags to a post
func (r *PostTagRepository) AddTagsToPost(postID uint, tags []domain.Tag) error {
	if len(tags) == 0 {
		return nil
	}
	// First remove existing tags for this post
	err := r.DB.Exec("DELETE FROM post_tags WHERE post_id = ?", postID).Error
	if err != nil {
		return err
	}
	// Then add new tags
	var post domain.Post
	post.ID = postID
	return r.DB.Model(&post).Association("Tags").Replace(tags)
}

// GetTagsByPostID returns tags for a post
func (r *PostTagRepository) GetTagsByPostID(postID uint) ([]domain.Tag, error) {
	var post domain.Post
	err := r.DB.Preload("Tags").First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return post.Tags, nil
}

// GetPostIDsByTagID returns post IDs for a tag
func (r *PostTagRepository) GetPostIDsByTagID(tagID uint) ([]uint, error) {
	var postIDs []uint
	err := r.DB.Raw(`
		SELECT post_id FROM post_tags WHERE tag_id = ?
	`, tagID).Scan(&postIDs).Error
	if err != nil {
		return nil, err
	}
	return postIDs, nil
}

// GetPostCountByTagID returns the number of posts for a tag
func (r *PostTagRepository) GetPostCountByTagID(tagID uint) (int64, error) {
	var count int64
	err := r.DB.Table("post_tags").Where("tag_id = ?", tagID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RemoveAllTagsFromPost removes all tags from a post
func (r *PostTagRepository) RemoveAllTagsFromPost(postID uint) error {
	return r.DB.Exec("DELETE FROM post_tags WHERE post_id = ?", postID).Error
}
