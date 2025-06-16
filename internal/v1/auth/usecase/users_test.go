package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/kadekchresna/payroll/config"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	"github.com/kadekchresna/payroll/internal/v1/auth/dto"
)

func Test_userUsecase_Create(t *testing.T) {
	now := time.Now()

	type args struct {
		ctx context.Context
		req dto.CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		// beforeFunc func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase
	}{
		{

			name: "CreateUsers-Success",
			args: args{
				ctx: context.Background(),
				req: dto.CreateUserRequest{
					Username: "username",
					Password: "password",
				},
			},
			// beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

			// 	return &userUsecase{
			// 		userRepo:     userRepo,
			// 		employeeRepo: employeeRepo,
			// 		time:         helper_time.NewTime(&now),
			// 		config:       config.Config{},
			// 	}
			// },
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// userRepo := new(repository_interface.)
			// employeeRepo := new(employee_repository_interface.MockIEmployeeRepository)

			uc := &userUsecase{
				// userRepo: userRepo,
				// employeeRepo: employeeRepo,
				time:   helper_time.NewTime(&now),
				config: config.Config{},
			}
			if err := uc.Create(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
