package domain

import "time"

// Topic 话题模型
type Topic struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Cover     string    `gorm:"size:255" json:"cover"`
	Description string `gorm:"size:500" json:"description"`
	PostCount int       `gorm:"not null;default:0" json:"post_count"`
	IsHot     bool      `gorm:"not null;default:false" json:"is_hot"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 关联
	Posts []Post `gorm:"many2many:post_topics;joinForeignKey:topic_id;joinReferences:post_id" json:"posts,omitempty"`
}

// TableName 指定表名
func (Topic) TableName() string {
	return "topics"
}
