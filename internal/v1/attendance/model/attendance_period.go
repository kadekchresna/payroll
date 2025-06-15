package model

import "time"

type AttendancePeriod struct {
	ID                 int       `json:"id"`
	PeriodStart        time.Time `json:"period_start"`
	PeriodEnd          time.Time `json:"period_end"`
	IsPayslipGenerated bool      `json:"is_payslip_generated"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedBy          int       `json:"created_by"`
	UpdatedBy          int       `json:"updated_by"`
}
