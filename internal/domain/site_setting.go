package domain

import "time"

// SiteSetting 站点配置项（key-value 结构）
type SiteSetting struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Key   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Group string `gorm:"type:varchar(50);index;not null;default:'general'" json:"group"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}
