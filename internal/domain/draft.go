package domain

import (
	"time"

	"gorm.io/gorm"
)

// Draft represents a saved post draft
type Draft struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID     uint   `gorm:"not null;index" json:"user_id"`
	Title      string `gorm:"type:varchar(200);not null;default:''" json:"title"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Category   string `gorm:"type:varchar(20);not null;default:'discussion'" json:"category"`
	PostType   string `gorm:"type:varchar(20);not null;default:'dynamic'" json:"post_type"`
	Visibility string `gorm:"type:varchar(20);not null;default:'public'" json:"visibility"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (Draft) TableName() string {
	return "drafts"
}
