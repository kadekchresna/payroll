package dao

import "time"

type PayslipDAO struct {
	ID                    int       `gorm:"column:id;primaryKey"`
	EmployeeID            int       `gorm:"column:employee_id"`
	TotalAttendanceDays   int       `gorm:"column:total_attendance_days"`
	TotalOvertimeHours    int       `gorm:"column:total_overtime_hours"`
	TotalAttendanceSalary float64   `gorm:"column:total_attendance_salary"`
	TotalOvertimeSalary   float64   `gorm:"column:total_overtime_salary"`
	TotalReimbursement    float64   `gorm:"column:total_reimbursement"`
	TotalTakeHomePay      float64   `gorm:"column:total_take_home_pay"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime:false"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy             int       `gorm:"column:created_by"`
	UpdatedBy             int       `gorm:"column:updated_by"`
	PeriodID              int       `gorm:"column:period_id"`
}

func (PayslipDAO) TableName() string {
	return "payslips"
}

type TotalTakeHomePayPerEmployee struct {
	EmployeeID       int     `gorm:"column:employee_id"`
	TotalTakeHomePay float64 `gorm:"column:total_take_home_pay"`
}
