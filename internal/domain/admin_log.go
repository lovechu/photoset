package domain

// AdminLog 操作日志
type AdminLog struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   uint   `gorm:"index" json:"admin_id"`
	AdminName string `gorm:"size:50" json:"admin_name"`
	Action    string `gorm:"size:50;not null;index" json:"action"`     // ban_user, unban_user, approve, reject, refund, role_change, delete_photoset, settings_update
	Target    string `gorm:"size:200" json:"target"`                    // 操作对象描述
	Detail    string `gorm:"type:text" json:"detail,omitempty"`         // 详细信息
	IP        string `gorm:"size:50" json:"ip,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

func (AdminLog) TableName() string {
	return "admin_logs"
}
