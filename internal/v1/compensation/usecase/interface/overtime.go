package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
)

type IOvertimeUsecase interface {
	CreateOvertime(ctx context.Context, req *dto.CreateOvertimeRequest) error
}
