package domain

import (
	"time"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusResolved ReportStatus = "resolved"
	ReportStatusRejected ReportStatus = "rejected"
)

// PostReport represents a report on a post or reply
type PostReport struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`

	PostID      *uint       `gorm:"index" json:"post_id"`   // Either post_id or reply_id should be set
	ReplyID     *uint       `gorm:"index" json:"reply_id"`  // Either post_id or reply_id should be set
	ReporterID  uint        `gorm:"not null;index" json:"reporter_id"`
	Reason      string      `gorm:"type:varchar(500);not null" json:"reason"`
	Status      string      `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	HandlerID   *uint       `json:"handler_id"`
	HandledAt   *time.Time  `json:"handled_at"`
	HandleNote  string      `gorm:"type:varchar(500)" json:"handle_note"`

	// Associations
	Post       *Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Reply      *PostReply `gorm:"foreignKey:ReplyID" json:"reply,omitempty"`
	Reporter   User       `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Handler    *User      `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
}

// TableName specifies the table name
func (PostReport) TableName() string {
	return "post_reports"
}

// Process processes the report (resolve or reject)
func (pr *PostReport) Process(handlerID uint, newStatus ReportStatus, note string) {
	now := time.Now()
	pr.Status = string(newStatus)
	pr.HandlerID = &handlerID
	pr.HandledAt = &now
	pr.HandleNote = note
}
