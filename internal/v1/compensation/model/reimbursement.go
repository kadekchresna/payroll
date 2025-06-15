package model

import "time"

type Reimbursement struct {
	ID          int       `json:"id"`
	EmployeeID  int       `json:"employee_id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
	PayslipID   int       `json:"payslip_id"`
}
