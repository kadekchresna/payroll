package dao

import (
	"time"
)

type OvertimeDAO struct {
	ID         int       `gorm:"primaryKey;column:id"`
	EmployeeID int       `gorm:"column:employee_id;not null;default:0"`
	Date       time.Time `gorm:"column:date;not null"`
	Hours      int       `gorm:"column:hours;not null;default:0"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime:false"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy  int       `gorm:"column:created_by;not null;default:0"`
	UpdatedBy  int       `gorm:"column:updated_by;not null;default:0"`
}

func (OvertimeDAO) TableName() string {
	return "overtimes"
}
