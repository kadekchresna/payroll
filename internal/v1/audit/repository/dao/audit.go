package dao

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         int            `gorm:"primaryKey;column:id"`
	TableNames string         `gorm:"column:table_name;not null"`
	Action     string         `gorm:"column:action;not null"`
	RecordID   int            `gorm:"column:record_id;not null"`
	OldData    datatypes.JSON `gorm:"column:old_data"`
	NewData    datatypes.JSON `gorm:"column:new_data"`
	ChangedBy  int            `gorm:"column:changed_by;not null"`
	ChangedAt  time.Time      `gorm:"column:changed_at;autoCreateTime:false"`
	IPAddress  string         `gorm:"column:ip_address;default:''"`
	RequestID  string         `gorm:"column:request_id;default:''"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
