package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
)

type IAttendancePeriodUsecase interface {
	CreateAttendancePeriod(ctx context.Context, req *dto.CreateAttendancePeriodRequest) error
}
