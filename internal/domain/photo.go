package domain

import "time"

// Photo 图片模型
type Photo struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	PhotoSetID uint       `gorm:"not null" json:"photoset_id"`
	URL        string     `gorm:"size:500;not null" json:"url"`
	SortOrder  int        `gorm:"default:0" json:"sort_order"`

	// 关联
	PhotoSet PhotoSet `gorm:"foreignKey:PhotoSetID" json:"photoset,omitempty"`
}

// TableName 指定表名
func (Photo) TableName() string {
	return "photos"
}
