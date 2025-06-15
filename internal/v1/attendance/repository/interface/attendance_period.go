package repository_interface

import (
	"context"
	"time"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
)

type IAttendancePeriodRepository interface {
	Create(ctx context.Context, p *model.AttendancePeriod) (int, error)
	GetByPeriodIntersect(ctx context.Context, periodStart time.Time, periodEnd time.Time) (*model.AttendancePeriod, error)
}
