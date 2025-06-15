package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/payslip/model"
)

type IPayslipRepository interface {
	Create(ctx context.Context, p *model.Payslip) (int, error)
	GetByID(ctx context.Context, id int) (*model.Payslip, error)
	GetTotalTakeHomePayPerEmployee(ctx context.Context) ([]model.TotalTakeHomePayPerEmployee, error)
	GetTotalTakeHomePayAllEmployees(ctx context.Context) (float64, error)
}
