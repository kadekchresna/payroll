package dto

import "time"

type CreateReimbursementRequest struct {
	EmployeeID  int
	UserID      int
	Amount      float64
	Description string
	Date        time.Time
}
