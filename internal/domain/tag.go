package domain

import "time"

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `gorm:"size:50;not null;uniqueIndex" json:"name"`

	// 关联
	PhotoSets []PhotoSet `gorm:"many2many:photoset_tags;joinForeignKey:photoset_id;joinReferences:tag_id" json:"photosets,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
