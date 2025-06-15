package model

import "time"

type Attendance struct {
	ID         int
	EmployeeID int
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CreatedBy  int
	UpdatedBy  int
}
