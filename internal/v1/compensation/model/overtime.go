package model

import "time"

type Overtime struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	Date       time.Time `json:"date"`
	Hours      int       `json:"hours"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  int       `json:"created_by"`
	UpdatedBy  int       `json:"updated_by"`
}

type EmployeeOvertimeSummary struct {
	EmployeeID int `json:"employee_id"`
	TotalHours int `json:"total_hours"`
}
