package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
)

type IAttendanceRepository interface {
	Create(ctx context.Context, a *model.Attendance) error
	GetByID(ctx context.Context, id int) (*model.Attendance, error)
	Update(ctx context.Context, a *model.Attendance) error
	Delete(ctx context.Context, id int) error
}
