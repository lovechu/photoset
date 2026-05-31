package repository

import (
	"photoset/internal/domain"

	"gorm.io/gorm"
)

// TopicRepository handles database operations for topics
type TopicRepository struct {
	DB *gorm.DB
}

// NewTopicRepository creates a new TopicRepository
func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{DB: db}
}

// Create creates a new topic
func (r *TopicRepository) Create(topic *domain.Topic) error {
	return r.DB.Create(topic).Error
}

// FindByID finds a topic by ID
func (r *TopicRepository) FindByID(id uint) (*domain.Topic, error) {
	var topic domain.Topic
	err := r.DB.First(&topic, id).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// FindByName finds a topic by name
func (r *TopicRepository) FindByName(name string) (*domain.Topic, error) {
	var topic domain.Topic
	err := r.DB.Where("name = ?", name).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// FindOrCreate finds a topic by name or creates it if it doesn't exist
func (r *TopicRepository) FindOrCreate(name string) (*domain.Topic, error) {
	var topic domain.Topic
	err := r.DB.Where("name = ?", name).First(&topic).Error
	if err == gorm.ErrRecordNotFound {
		topic = domain.Topic{Name: name, Status: "active"}
		if err := r.DB.Create(&topic).Error; err != nil {
			return nil, err
		}
		return &topic, nil
	}
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// FindOrCreateBatch finds or creates multiple topics by name
func (r *TopicRepository) FindOrCreateBatch(names []string) ([]domain.Topic, error) {
	var topics []domain.Topic
	for _, name := range names {
		topic, err := r.FindOrCreate(name)
		if err != nil {
			return nil, err
		}
		topics = append(topics, *topic)
	}
	return topics, nil
}

// GetAll returns all topics
func (r *TopicRepository) GetAll() ([]domain.Topic, error) {
	var topics []domain.Topic
	err := r.DB.Where("status = ?", "active").Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

// GetHotTopics returns hot topics sorted by post count
func (r *TopicRepository) GetHotTopics(limit int) ([]domain.Topic, error) {
	var topics []domain.Topic
	err := r.DB.
		Where("status = ?", "active").
		Order("post_count DESC").
		Limit(limit).
		Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

// SearchByName searches topics by name (partial match)
func (r *TopicRepository) SearchByName(keyword string, limit int) ([]domain.Topic, error) {
	var topics []domain.Topic
	err := r.DB.Where("name LIKE ? AND status = ?", "%"+keyword+"%", "active").
		Limit(limit).
		Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

// Update updates a topic
func (r *TopicRepository) Update(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.Topic{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a topic by ID
func (r *TopicRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Topic{}, id).Error
}

// IncrementPostCount increments the post count for a topic
func (r *TopicRepository) IncrementPostCount(id uint) error {
	return r.DB.Model(&domain.Topic{}).Where("id = ?", id).Update("post_count", gorm.Expr("post_count + 1")).Error
}

// DecrementPostCount decrements the post count for a topic
func (r *TopicRepository) DecrementPostCount(id uint) error {
	return r.DB.Model(&domain.Topic{}).Where("id = ?", id).Update("post_count", gorm.Expr("GREATEST(post_count - 1, 0)")).Error
}
