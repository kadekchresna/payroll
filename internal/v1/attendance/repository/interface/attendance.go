package repository_interface

import (
	"context"
	"time"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
)

type IAttendanceRepository interface {
	Create(ctx context.Context, a *model.Attendance) (int, error)
	GetByID(ctx context.Context, id int) (*model.Attendance, error)
	GetByDateAndEmployeeID(ctx context.Context, employeeID int, date time.Time) (*model.Attendance, error)
	GetEmployeeCountByDateRange(ctx context.Context, periodStart time.Time, periodEnd time.Time) ([]*model.EmployeeAttendanceCount, error)
}
