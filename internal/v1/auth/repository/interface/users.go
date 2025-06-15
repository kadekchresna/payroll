package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/auth/model"
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}
