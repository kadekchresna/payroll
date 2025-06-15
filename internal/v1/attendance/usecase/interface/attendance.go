package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
)

type IAttendanceUsecase interface {
	CreateAttendance(ctx context.Context, a *model.Attendance) error
	GetAttendance(ctx context.Context, id int) (*model.Attendance, error)
	UpdateAttendance(ctx context.Context, a *model.Attendance) error
	DeleteAttendance(ctx context.Context, id int) error
}
