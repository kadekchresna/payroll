package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/employee/model"
)

type IEmployeeRepository interface {
	Create(ctx context.Context, e *model.Employee) error
	GetByID(ctx context.Context, id int) (*model.Employee, error)
	Update(ctx context.Context, e *model.Employee) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]model.Employee, error)
	GetByUserID(ctx context.Context, userID int) (*model.Employee, error)
}
