package dao

import "time"

type AttendancePeriodDAO struct {
	ID                 int       `gorm:"primaryKey;column:id"`
	PeriodStart        time.Time `gorm:"column:period_start;default:now()"`
	PeriodEnd          time.Time `gorm:"column:period_end;default:now()"`
	IsPayslipGenerated *bool     `gorm:"column:is_payslip_generated;default:false"`
	CreatedAt          time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy          int       `gorm:"column:created_by;default:0"`
	UpdatedBy          int       `gorm:"column:updated_by;default:0"`
}

func (AttendancePeriodDAO) TableName() string {
	return "attendances_period"
}
