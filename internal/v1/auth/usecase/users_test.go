package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	"github.com/kadekchresna/payroll/internal/v1/auth/dto"
	auth_model "github.com/kadekchresna/payroll/internal/v1/auth/model"

	user_repository_interface "github.com/kadekchresna/payroll/internal/v1/auth/repository/interface/mocks"
	employee_model "github.com/kadekchresna/payroll/internal/v1/employee/model"
	employee_repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_userUsecase_Create(t *testing.T) {
	now := time.Now()

	type args struct {
		ctx context.Context
		req dto.CreateUserRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase
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
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config:       config.Config{},
				}
			},
			wantErr: false,
		},
		{

			name: "CreateUsers-Failed",
			args: args{
				ctx: context.Background(),
				req: dto.CreateUserRequest{
					Username: "username",
					Password: "password",
				},
			},
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.New("FATAL ERROR"))

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config:       config.Config{},
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepo := new(user_repository_interface.MockIUserRepository)
			employeeRepo := new(employee_repository_interface.MockIEmployeeRepository)

			uc := tt.beforeFunc(userRepo, employeeRepo)
			if err := uc.Create(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUsecase_Login(t *testing.T) {
	now := time.Now()

	u := auth_model.User{
		ID:       1,
		Username: "username",
		Password: "xkvi2+l/y17mB3FVxULKaiG4KL/L01cM23Synv8raIE=",
		Salt:     "ENpjbttXJDkGPjV4ablfmA==",
	}

	employee := employee_model.Employee{
		ID:       1,
		FullName: "Fullname",
		UserID:   u.ID,
	}

	accessToken, refreshToken, _ := jwt.GenerateToken("secret", u.ID, u.Role, employee.ID, employee.FullName)

	type args struct {
		ctx context.Context
		req dto.LoginUserRequest
	}
	tests := []struct {
		name       string
		args       args
		want       *dto.LoginUserResponse
		wantErr    bool
		beforeFunc func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase
	}{
		{
			name: "Login-Success",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret123",
				},
			},
			want: &dto.LoginUserResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
			wantErr: false,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(&u, nil).Once()
				employeeRepo.EXPECT().GetByUserID(mock.Anything, u.ID).Return(&employee, nil).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-SuccessWihtoutEmployee",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret123",
				},
			},
			want: func() *dto.LoginUserResponse {
				accessTokenWtEmployee, refreshTokenWtEmployee, _ := jwt.GenerateToken("secret", u.ID, u.Role, 0, "")
				return &dto.LoginUserResponse{
					AccessToken:  accessTokenWtEmployee,
					RefreshToken: refreshTokenWtEmployee,
				}
			}(),
			wantErr: false,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(&u, nil).Once()
				employeeRepo.EXPECT().GetByUserID(mock.Anything, u.ID).Return(nil, nil).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedErrorQueryEmployeeGetByUserID",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret123",
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(&u, nil).Once()
				employeeRepo.EXPECT().GetByUserID(mock.Anything, u.ID).Return(nil, errors.New("FATAL ERROR")).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedPasswordNotMatch",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret1234",
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(&u, nil).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedErrorQueryUserGetByUsername",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret1234",
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(nil, errors.New("FATAL ERROR")).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedUserNotFound",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "secret1234",
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				userRepo.EXPECT().GetByUsername(mock.Anything, u.Username).Return(nil, nil).Once()

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedPasswordEmpty",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: u.Username,
					Password: "",
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
		{
			name: "Login-FailedUsernameEmpty",
			args: args{
				ctx: context.Background(),
				req: dto.LoginUserRequest{
					Username: "",
					Password: u.Password,
				},
			},
			want: func() *dto.LoginUserResponse {
				return nil
			}(),
			wantErr: true,
			beforeFunc: func(userRepo *user_repository_interface.MockIUserRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository) *userUsecase {

				return &userUsecase{
					userRepo:     userRepo,
					employeeRepo: employeeRepo,
					time:         helper_time.NewTime(&now),
					config: config.Config{
						AppJWTSecret: "secret",
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepo := new(user_repository_interface.MockIUserRepository)
			employeeRepo := new(employee_repository_interface.MockIEmployeeRepository)

			uc := tt.beforeFunc(userRepo, employeeRepo)
			got, err := uc.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUsecase.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
