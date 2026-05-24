package domain

import (
	"strings"
	"time"
)

// SensitiveWord represents a sensitive word for content filtering
type SensitiveWord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Word       string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"word"`
	Replacement string   `gorm:"type:varchar(100);not null;default:'***'" json:"replacement"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name
func (SensitiveWord) TableName() string {
	return "sensitive_words"
}

// Match checks if the given text contains this sensitive word (case-insensitive)
func (sw *SensitiveWord) Match(text string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(sw.Word))
}

// Replace replaces sensitive words in the given text with the replacement
func (sw *SensitiveWord) Replace(text string) string {
	return strings.ReplaceAll(strings.ToLower(text), strings.ToLower(sw.Word), sw.Replacement)
}
