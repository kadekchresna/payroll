package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/employee/model"
)

type IEmployeeRepository interface {
	GetByID(ctx context.Context, id int) (*model.Employee, error)
	GetByUserID(ctx context.Context, userID int) (*model.Employee, error)
}
