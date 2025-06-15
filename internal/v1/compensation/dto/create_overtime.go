package dto

import "time"

type CreateOvertimeRequest struct {
	EmployeeID int
	UserID     int
	Date       time.Time
	Hours      int
}
