package dto

import "time"

type CreateAttendancePeriodRequest struct {
	PeriodStart time.Time
	PeriodEnd   time.Time
	UserID      int
}
