package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
)

type IReimbursementUsecase interface {
	CreateReimbursement(ctx context.Context, req *dto.CreateReimbursementRequest) error
}
