package dto

import "time"

type CreateAttendanceRequest struct {
	EmployeeID int
	UserID     int
	Date       time.Time
}
