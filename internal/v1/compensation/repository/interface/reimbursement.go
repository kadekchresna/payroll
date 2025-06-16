package repository_interface

import (
	"context"
	"time"

	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
)

type IReimbursementRepository interface {
	Create(ctx context.Context, m *model.Reimbursement) (int, error)
	SumReimbursementsByID(ctx context.Context, id []int) ([]*model.EmployeeReimbursementSummary, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]model.Reimbursement, error)
	Update(ctx context.Context, m *model.Reimbursement, ids []int) error
	GetByIDs(ctx context.Context, id []int) ([]model.Reimbursement, error)
	GetByPayslipID(ctx context.Context, payslipID int, employeeID int) ([]model.Reimbursement, error)
}
