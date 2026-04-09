package domain

import "time"

// MembershipPlan 会员套餐模型
type MembershipPlan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `gorm:"size:50;not null" json:"name"`                                     // 套餐名：月度会员、年度会员
	Duration    int       `gorm:"not null;comment:会员天数" json:"duration"`                          // 30 或 365
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`                          // 价格
	Description string    `gorm:"type:text" json:"description"`                                      // 套餐说明
	Status      int       `gorm:"type:tinyint;default:1;comment:1-on,0-off" json:"status"`           // 上架/下架
}

func (MembershipPlan) TableName() string {
	return "memberships"
}
