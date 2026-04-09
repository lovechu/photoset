package domain

import "time"

// Order 订单模型
type Order struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	OrderNo       string     `gorm:"size:64;uniqueIndex;not null" json:"order_no"`                            // 订单号
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	Type          string     `gorm:"size:20;not null" json:"type"`                                           // membership / single
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`                               // 金额
	Status        string     `gorm:"size:20;default:pending" json:"status"`                                  // pending/paid/cancelled/refunded
	MembershipID  *uint      `json:"membership_id"`                                                           // 会员套餐 ID（type=membership 时有值）
	PhotoSetID    *uint      `json:"photoset_id"`                                                             // 套图 ID（type=single 时有值）
	PaidAt        *time.Time `json:"paid_at"`                                                                 // 支付时间
	ExpireSeconds int        `gorm:"default:1800;comment:支付过期秒数" json:"-"`                                // 订单过期时间（默认30分钟）

	// 关联
	User       User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Membership *MembershipPlan `gorm:"foreignKey:MembershipID" json:"membership,omitempty"`
	PhotoSet   *PhotoSet       `gorm:"foreignKey:PhotoSetID" json:"photoset,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
