package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// PostTopicRepository handles database operations for post-topic relationships
type PostTopicRepository struct {
	DB *gorm.DB
}

// NewPostTopicRepository creates a new PostTopicRepository
func NewPostTopicRepository(db *gorm.DB) *PostTopicRepository {
	return &PostTopicRepository{DB: db}
}

// AddTopicsToPost adds topics to a post
func (r *PostTopicRepository) AddTopicsToPost(postID uint, topics []domain.Topic) error {
	if len(topics) == 0 {
		return nil
	}
	// First remove existing topics for this post
	err := r.DB.Exec("DELETE FROM post_topics WHERE post_id = ?", postID).Error
	if err != nil {
		return err
	}
	// Then add new topics
	var post domain.Post
	post.ID = postID
	return r.DB.Model(&post).Association("Topics").Replace(topics)
}

// GetTopicsByPostID returns topics for a post
func (r *PostTopicRepository) GetTopicsByPostID(postID uint) ([]domain.Topic, error) {
	var post domain.Post
	err := r.DB.Preload("Topics").First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return post.Topics, nil
}

// GetPostIDsByTopicID returns post IDs for a topic
func (r *PostTopicRepository) GetPostIDsByTopicID(topicID uint) ([]uint, error) {
	var postIDs []uint
	err := r.DB.Raw(`
		SELECT post_id FROM post_topics WHERE topic_id = ?
	`, topicID).Scan(&postIDs).Error
	if err != nil {
		return nil, err
	}
	return postIDs, nil
}

// GetPostCountByTopicID returns the number of posts for a topic
func (r *PostTopicRepository) GetPostCountByTopicID(topicID uint) (int64, error) {
	var count int64
	err := r.DB.Table("post_topics").Where("topic_id = ?", topicID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// RemoveAllTopicsFromPost removes all topics from a post
func (r *PostTopicRepository) RemoveAllTopicsFromPost(postID uint) error {
	return r.DB.Exec("DELETE FROM post_topics WHERE post_id = ?", postID).Error
}
