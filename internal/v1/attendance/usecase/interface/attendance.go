package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
)

type IAttendanceUsecase interface {
	CreateAttendance(ctx context.Context, req *dto.CreateAttendanceRequest) error
}
