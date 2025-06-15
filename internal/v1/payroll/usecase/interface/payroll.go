package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/payroll/dto"
)

type IPayrollUsecase interface {
	CreatePayroll(ctx context.Context, req dto.CreatePayrollRequest) error
	GetPayrollByID(ctx context.Context, req *dto.GetEmployeePayrollRequest) (*dto.GetEmployeePayrollResponse, error)
	GetPayrollSummary(ctx context.Context) (*dto.GetPayrollSummary, error)
}
