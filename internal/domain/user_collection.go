package domain

import "time"

// UserCollection 用户合集
type UserCollection struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index:idx_user_collection" json:"user_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	CoverURL    string    `gorm:"size:500" json:"cover_url"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	ItemCount   int       `gorm:"default:0" json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	User  User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items []CollectionItem  `gorm:"foreignKey:CollectionID" json:"items,omitempty"`
}

func (UserCollection) TableName() string {
	return "user_collections"
}

// CollectionItem 合集项目
type CollectionItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CollectionID uint      `gorm:"not null;index:idx_collection_photoset" json:"collection_id"`
	PhotoSetID   uint      `gorm:"not null;index:idx_collection_photoset" json:"photoset_id"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联
	PhotoSet PhotoSet `gorm:"foreignKey:PhotoSetID" json:"photoset,omitempty"`
}

func (CollectionItem) TableName() string {
	return "collection_items"
}
