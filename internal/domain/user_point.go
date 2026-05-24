package domain

import (
	"time"
)

// UserPoint represents user points and level
type UserPoint struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	Points    int       `gorm:"not null;default:0" json:"points"`
	Level     int       `gorm:"not null;default:1" json:"level"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// UserPointLog represents a log of point changes
type UserPointLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	Points    int       `gorm:"not null" json:"points"`
	Action    string    `gorm:"type:varchar(50);not null" json:"action"`
	RelatedID *uint     `gorm:"column:related_id" json:"related_id,omitempty"` // Post ID, Reply ID, etc.
}

// TableName specifies the table name for UserPointLog
func (UserPointLog) TableName() string {
	return "user_point_logs"
}

// TableName specifies the table name
func (UserPoint) TableName() string {
	return "user_points"
}

// LevelThresholds defines the points required for each level
var LevelThresholds = []int{
	0,      // Level 1: 0+ points
	100,    // Level 2: 100+ points
	500,    // Level 3: 500+ points
	2000,   // Level 4: 2000+ points
	5000,   // Level 5: 5000+ points
	10000,  // Level 6: 10000+ points
	20000,  // Level 7: 20000+ points
	30000,  // Level 8: 30000+ points
	40000,  // Level 9: 40000+ points
	50000,  // Level 10: 50000+ points
}

// LevelNames defines the names for each level
var LevelNames = []string{
	"",           // placeholder for index 0
	"Newbie",     // Level 1
	"Active Member",   // Level 2
	"Senior Member",   // Level 3
	"Gold Member",     // Level 4
	"Diamond Member",  // Level 5
	"Supreme Member",  // Level 6
	"Glory L7",       // Level 7
	"Glory L8",       // Level 8
	"Glory L9",       // Level 9
	"Glory L10",      // Level 10
}

// CalculateLevel calculates the level based on points
func CalculateLevel(points int) int {
	for i := len(LevelThresholds) - 1; i >= 0; i-- {
		if points >= LevelThresholds[i] {
			return i + 1 // Level starts from 1
		}
	}
	return 1 // Default to level 1
}

// GetLevelName returns the name of the given level
func GetLevelName(level int) string {
	if level < 1 || level >= len(LevelNames) {
		return "Unknown"
	}
	return LevelNames[level]
}

// AddPoints adds points and updates the level
func (up *UserPoint) AddPoints(points int) {
	up.Points += points
	if up.Points < -9999 {
		up.Points = -9999
	}
	up.Level = CalculateLevel(up.Points)
	up.UpdatedAt = time.Now()
}

// GetNextLevelPoints returns the points needed for the next level
func (up *UserPoint) GetNextLevelPoints() int {
	currentLevel := up.Level
	if currentLevel >= len(LevelThresholds) {
		return 0 // Max level reached
	}
	return LevelThresholds[currentLevel]
}
