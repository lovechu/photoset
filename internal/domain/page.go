package domain

type Page struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Slug      string `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Title     string `gorm:"size:200;not null" json:"title"`
	ContentMD string `gorm:"type:text" json:"content_md"`
	Content   string `gorm:"-" json:"content,omitempty"` // 仅用于 API 返回的渲染后 HTML，不存数据库
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Status    string `gorm:"size:20;default:'published';index" json:"status"` // published/draft
	CreatedAt int64  `gorm:"type:bigint" json:"created_at"`
	UpdatedAt int64  `gorm:"type:bigint" json:"updated_at"`
}

const (
	PageStatusPublished = "published"
	PageStatusDraft     = "draft"
)