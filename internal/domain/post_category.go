package domain

import "time"

// CommunityCategory 社区帖子分类模型（独立于套图的 Category）
type CommunityCategory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"size:64;not null;uniqueIndex" json:"key"`          // 标识，如 "discussion"
	Name        string    `gorm:"size:128;not null" json:"name"`                     // 显示名称，如 "讨论"
	Description string    `gorm:"size:256" json:"description,omitempty"`             // 描述
	Color       string    `gorm:"size:7;default:'#409EFF'" json:"color"`             // hex 颜色
	Icon        string    `gorm:"size:64;default:''" json:"icon,omitempty"`          // 图标
	SortOrder   int       `gorm:"default:0" json:"sort_order"`                       // 排序号
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the database table name for CommunityCategory
func (CommunityCategory) TableName() string {
	return "post_categories"
}
