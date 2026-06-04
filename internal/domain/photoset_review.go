package domain

import "time"

// PhotoSetReview 套图评价
type PhotoSetReview struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index:idx_user_photoset" json:"user_id"`
	PhotoSetID  uint      `gorm:"not null;index:idx_user_photoset" json:"photoset_id"`
	Rating      int       `gorm:"not null" json:"rating"`       // 1-5 星评分
	Tags        string    `gorm:"type:varchar(500);default:''" json:"tags"` // 评价标签，逗号分隔
	Content     string    `gorm:"type:text" json:"content"`     // 评价内容
	IsAnonymous bool      `gorm:"default:false" json:"is_anonymous"` // 是否匿名
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PhotoSet *PhotoSet `gorm:"foreignKey:PhotoSetID" json:"photoset,omitempty"`
}

func (PhotoSetReview) TableName() string {
	return "photoset_reviews"
}

// ReviewSummary 评价汇总
type ReviewSummary struct {
	TotalCount    int64   `json:"total_count"`
	AverageRating float64 `json:"average_rating"`
	Rating5Count  int64   `json:"rating_5_count"`
	Rating4Count  int64   `json:"rating_4_count"`
	Rating3Count  int64   `json:"rating_3_count"`
	Rating2Count  int64   `json:"rating_2_count"`
	Rating1Count  int64   `json:"rating_1_count"`
	TopTags       []TagCount `json:"top_tags"`
}

// TagCount 标签统计
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}
