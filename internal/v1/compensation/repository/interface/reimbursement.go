package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
)

type IReimbursementRepository interface {
	Create(ctx context.Context, m *model.Reimbursement) (int, error)
}
