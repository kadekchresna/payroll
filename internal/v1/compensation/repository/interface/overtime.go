package repository_interface

import (
	"context"
	"time"

	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
)

type IOvertimeRepository interface {
	Create(ctx context.Context, ot *model.Overtime) (int, error)
	GetByDateAndEmployeeID(ctx context.Context, employeeID int, date time.Time) ([]model.Overtime, error)
}
