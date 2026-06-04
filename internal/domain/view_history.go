package domain

import "time"

// ViewHistory 套图浏览历史记录
type ViewHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      uint      `gorm:"not null;index:idx_user_viewed,priority:1" json:"user_id"`
	PhotosetID  uint      `gorm:"not null;index:idx_user_viewed,priority:2" json:"photoset_id"`
	ViewedAt    time.Time `gorm:"not null;index:idx_user_viewed,priority:3" json:"viewed_at"`

	// 冗余字段（用于列表展示，避免 JOIN 查询）
	Title       string    `gorm:"type:varchar(200);not null" json:"title"`
	CoverImage  string    `gorm:"type:varchar(500);default:''" json:"cover_image"`
	PhotoCount  int       `gorm:"default:0" json:"photo_count"`
	IsVip       bool      `gorm:"default:false" json:"is_vip"`
	Price       float64   `gorm:"type:decimal(10,2);default:0" json:"price"`

	// 关联
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	PhotoSet  PhotoSet  `gorm:"foreignKey:PhotosetID" json:"-"`
}

// TableName 指定表名
func (ViewHistory) TableName() string {
	return "view_histories"
}
