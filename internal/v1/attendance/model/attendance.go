package model

import "time"

type Attendance struct {
	ID           int        `json:"id"`
	EmployeeID   int        `json:"employee_id"`
	Date         time.Time  `json:"date"`
	CreatedAt    time.Time  `json:"created_at"`
	CheckedInAt  time.Time  `json:"checked_in_at"`
	CheckedOutAt *time.Time `json:"checked_out_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    int        `json:"created_by"`
	UpdatedBy    int        `json:"updated_by"`
}
