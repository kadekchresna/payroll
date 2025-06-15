package dto

import "time"

type GetEmployeePayrollRequest struct {
	EmployeeID int
	PayslipID  int
}

type Reimbursement struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

type GetEmployeePayrollResponse struct {
	TotalWorkingDaysInAttendancePeriod int             `json:"total_working_days_in_period"`
	TotalSalaryPerMonth                float64         `json:"total_salary_per_month"`
	TotalAttendanceDays                int             `json:"total_attendance_days"`
	TotalOvertimeHours                 int             `json:"total_overtime_hours"`
	TotalAttendanceSalary              float64         `json:"total_attendance_salary"`
	TotalOvertimeSalary                float64         `json:"total_overtime_salary"`
	TotalReimbursement                 float64         `json:"total_reimbursement"`
	TotalTakeHomePay                   float64         `json:"total_take_home_pay"`
	PeriodStart                        time.Time       `json:"period_start"`
	PeriodEnd                          time.Time       `json:"period_end"`
	Reimbursement                      []Reimbursement `json:"reimbursements"`
}
