package model

import "time"

type Payslip struct {
	ID                    int       `json:"id"`
	EmployeeID            int       `json:"employee_id"`
	TotalAttendanceDays   int       `json:"total_attendance_days"`
	TotalOvertimeHours    int       `json:"total_overtime_hours"`
	TotalAttendanceSalary float64   `json:"total_attendance_salary"`
	TotalOvertimeSalary   float64   `json:"total_overtime_salary"`
	TotalReimbursement    float64   `json:"total_reimbursement"`
	TotalTakeHomePay      float64   `json:"total_take_home_pay"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	CreatedBy             int       `json:"created_by"`
	UpdatedBy             int       `json:"updated_by"`
	PeriodID              int       `json:"period_id"`
}

type TotalTakeHomePayPerEmployee struct {
	EmployeeID       int     `json:"employee_id"`
	TotalTakeHomePay float64 `json:"total_take_home_pay"`
}
