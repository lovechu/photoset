package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleGuest   UserRole = "guest"
	RoleUser    UserRole = "user"
	RoleMember  UserRole = "member"
	RoleCreator UserRole = "creator"
	RoleVIP     UserRole = "vip"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Nickname     string    `gorm:"type:varchar(50);not null" json:"nickname"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole  `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	Status       int       `gorm:"type:tinyint;default:1;comment:1-active,0-inactive" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	MembershipExpires *time.Time `gorm:"index" json:"membership_expires"` // 会员过期时间，nil 表示非会员

	Avatar       string    `gorm:"type:varchar(255);default:''" json:"avatar"`
	Bio          string    `gorm:"type:varchar(500);default:''" json:"bio"`
	IPLocation   string    `gorm:"type:varchar(50);default:''" json:"ip_location"`
	Level        int       `gorm:"type:tinyint;default:1" json:"level"`
	Exp          int       `gorm:"type:int;default:0" json:"exp"`
	CircleCount  int       `gorm:"type:int;default:0" json:"circle_count"`
	FollowingCount int     `gorm:"type:int;default:0" json:"following_count"`
	FollowerCount  int     `gorm:"type:int;default:0" json:"follower_count"`
	LikeCount    int       `gorm:"type:int;default:0" json:"like_count"`
}

func (User) TableName() string {
	return "users"
}
