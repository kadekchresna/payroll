package dao

import (
	"time"
)

type ReimbursementDAO struct {
	ID          int       `gorm:"primaryKey;column:id"`
	EmployeeID  int       `gorm:"column:employee_id;not null;default:0"`
	Date        time.Time `gorm:"column:date;not null"`
	Amount      float64   `gorm:"column:amount;not null;default:0"`
	Description string    `gorm:"column:description;not null;default:''"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy   int       `gorm:"column:created_by;not null;default:0"`
	UpdatedBy   int       `gorm:"column:updated_by;not null;default:0"`
	PayslipID   int       `gorm:"column:payslip_id;not null;default:0"`
}

func (ReimbursementDAO) TableName() string {
	return "reimbursements"
}
