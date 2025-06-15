package dao

import (
	"database/sql"
	"time"
)

type AttendanceDAO struct {
	ID           int          `gorm:"column:id;primaryKey"`
	EmployeeID   int          `gorm:"column:employee_id;not null;index:idx_employee_date,unique"`
	Date         time.Time    `gorm:"column:date;not null;index:idx_employee_date,unique"`
	CheckedInAt  time.Time    `gorm:"column:checked_in_at;not null;"`
	CheckedOutAt sql.NullTime `gorm:"column:checked_out_at"`
	CreatedAt    time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy    int          `gorm:"column:created_by;default:0"`
	UpdatedBy    int          `gorm:"column:updated_by;default:0"`
}

func (AttendanceDAO) TableName() string {
	return "attendances"
}
