package domain

import "time"

// Category 分类模型（类似 Tag 但用于单选分类）
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"size:50;not null;uniqueIndex" json:"name"`     // 显示名称，如 "自然风光"
	Slug        string    `gorm:"size:50;not null;uniqueIndex" json:"slug"`     // URL标识，如 "nature"
	Description string    `gorm:"size:200" json:"description,omitempty"`        // 描述（可选）
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                  // 排序权重，越大越前

	// 注意：这里不能直接包含PhotoSet切片，会导致循环引用
	// 如果需要关联查询，应该使用单独的查询语句
}

func (Category) TableName() string {
	return "categories"
}
