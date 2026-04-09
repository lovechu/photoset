package domain

import "time"

// Favorite 收藏模型
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_user_photoset" json:"user_id"`
	PhotoSetID uint     `gorm:"not null;uniqueIndex:uk_user_photoset;column:photoset_id" json:"photoset_id"`
	PhotoSet  PhotoSet  `gorm:"foreignKey:PhotoSetID" json:"photoset,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}
