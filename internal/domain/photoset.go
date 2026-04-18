package domain

import "time"

// PhotoSet 套图模型
type PhotoSet struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Cover       string    `gorm:"size:500;not null" json:"cover"`
	Description string    `gorm:"type:text" json:"description"`
	IsFree      int8      `gorm:"default:1;comment:1-free,0-paid" json:"is_free"`
	Price       float64   `gorm:"type:decimal(10,2);default:0" json:"price"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Status      string    `gorm:"size:20;default:draft;comment:draft,published,pending" json:"status"`
	Category    string    `gorm:"size:50;default:'';index" json:"category"` // 分类 slug，如 'nature'

	// 非数据库字段（忽略数据库的 photo_count 列，手动计算）
	IsFavorited bool `gorm:"-" json:"is_favorited"` // 是否已收藏
	PhotoCount  int  `gorm:"-" json:"photo_count" gorm:"-"`  // 图片数量（忽略数据库列）

	// 关联
	User   User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Photos []Photo  `gorm:"foreignKey:photoset_id;References:ID" json:"photos,omitempty"`
	Tags   []Tag    `gorm:"many2many:photoset_tags;joinForeignKey:photoset_id;joinReferences:tag_id" json:"tags,omitempty"`
}

// TableName 指定表名
func (PhotoSet) TableName() string {
	return "photosets"
}
