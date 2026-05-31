package domain

import (
	"time"

	"gorm.io/gorm"
)

// PointsMallItem defines an item in the points mall
type PointsMallItem struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Image       string `gorm:"type:varchar(255);default:''" json:"image"`
	Category    string `gorm:"type:varchar(50);not null;index" json:"category"` // badge, title, privilege, virtual_gift
	
	// Points cost
	PointsCost int `gorm:"not null;default:0" json:"points_cost"`
	
	// Item type and value
	ItemType  string `gorm:"type:varchar(50);not null" json:"item_type"` // badge, title, vip_days, custom
	ItemValue string `gorm:"type:varchar(255);not null" json:"item_value"`
	
	// Stock management
	TotalStock  int  `gorm:"default:-1" json:"total_stock"` // -1 means unlimited
	UsedStock   int  `gorm:"default:0" json:"used_stock"`
	IsUnlimited bool `gorm:"default:true" json:"is_unlimited"`
	
	// Requirements
	MinLevel int `gorm:"default:1" json:"min_level"`
	
	// Status
	IsActive  bool `gorm:"default:true" json:"is_active"`
	SortOrder int  `gorm:"default:0" json:"sort_order"`
}

// TableName specifies the table name
func (PointsMallItem) TableName() string {
	return "points_mall_items"
}

// UserPointsExchange represents a user's points exchange record
type UserPointsExchange struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	
	UserID    uint `gorm:"not null;index" json:"user_id"`
	ItemID    uint `gorm:"not null" json:"item_id"`
	Points    int  `gorm:"not null" json:"points"`
	
	// Exchange details
	ItemName  string `gorm:"type:varchar(100)" json:"item_name"`
	ItemType  string `gorm:"type:varchar(50)" json:"item_type"`
	ItemValue string `gorm:"type:varchar(255)" json:"item_value"`
	
	// Status
	Status    string `gorm:"type:varchar(20);default:'completed'" json:"status"` // completed, refunded
	
	// Associations
	User      User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Item      PointsMallItem   `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

// TableName specifies the table name
func (UserPointsExchange) TableName() string {
	return "user_points_exchanges"
}

// GetRemainingStock returns remaining stock
func (item *PointsMallItem) GetRemainingStock() int {
	if item.IsUnlimited {
		return -1
	}
	remaining := item.TotalStock - item.UsedStock
	if remaining < 0 {
		return 0
	}
	return remaining
}

// CanExchange checks if user can exchange this item
func (item *PointsMallItem) CanExchange(userLevel int, userPoints int) bool {
	if !item.IsActive {
		return false
	}
	if item.MinLevel > userLevel {
		return false
	}
	if item.PointsCost > userPoints {
		return false
	}
	if !item.IsUnlimited && item.GetRemainingStock() <= 0 {
		return false
	}
	return true
}

// DefaultPointsMallItems returns default items for the points mall
func DefaultPointsMallItems() []PointsMallItem {
	return []PointsMallItem{
		// Badges
		{
			Name:        "专属徽章-新手上路",
			Description: "新手专属徽章，展示你的社区起点",
			Category:    "badge",
			PointsCost:  100,
			ItemType:    "badge",
			ItemValue:   "badge_newbie",
			MinLevel:    1,
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Name:        "专属徽章-活跃达人",
			Description: "活跃达人徽章，展示你的社区活力",
			Category:    "badge",
			PointsCost:  500,
			ItemType:    "badge",
			ItemValue:   "badge_active",
			MinLevel:    2,
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Name:        "专属徽章-创作大师",
			Description: "创作大师徽章，展示你的创作才华",
			Category:    "badge",
			PointsCost:  2000,
			ItemType:    "badge",
			ItemValue:   "badge_creator",
			MinLevel:    4,
			IsActive:    true,
			SortOrder:   3,
		},
		
		// Titles
		{
			Name:        "称号-社区之星",
			Description: "闪耀的社区之星称号",
			Category:    "title",
			PointsCost:  300,
			ItemType:    "title",
			ItemValue:   "社区之星",
			MinLevel:    2,
			IsActive:    true,
			SortOrder:   10,
		},
		{
			Name:        "称号-创作先锋",
			Description: "勇于创作的先锋称号",
			Category:    "title",
			PointsCost:  800,
			ItemType:    "title",
			ItemValue:   "创作先锋",
			MinLevel:    3,
			IsActive:    true,
			SortOrder:   11,
		},
		{
			Name:        "称号-意见领袖",
			Description: "社区意见领袖称号",
			Category:    "title",
			PointsCost:  1500,
			ItemType:    "title",
			ItemValue:   "意见领袖",
			MinLevel:    4,
			IsActive:    true,
			SortOrder:   12,
		},
		
		// Privileges
		{
			Name:        "VIP体验卡-7天",
			Description: "7天VIP体验，享受会员特权",
			Category:    "privilege",
			PointsCost:  1000,
			ItemType:    "vip_days",
			ItemValue:   "7",
			MinLevel:    3,
			IsActive:    true,
			SortOrder:   20,
		},
		{
			Name:        "VIP体验卡-30天",
			Description: "30天VIP体验，享受会员特权",
			Category:    "privilege",
			PointsCost:  3000,
			ItemType:    "vip_days",
			ItemValue:   "30",
			MinLevel:    5,
			IsActive:    true,
			SortOrder:   21,
		},
		{
			Name:        "自定义头像框",
			Description: "解锁专属头像框装饰",
			Category:    "privilege",
			PointsCost:  2000,
			ItemType:    "custom",
			ItemValue:   "custom_avatar_frame",
			MinLevel:    4,
			IsActive:    true,
			SortOrder:   22,
		},
	}
}
