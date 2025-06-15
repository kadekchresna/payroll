package usecase_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/auth/dto"
)

type IUserUsecase interface {
	Create(ctx context.Context, req dto.CreateUserRequest) error
	Login(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error)
}
