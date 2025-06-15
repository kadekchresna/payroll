package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/helper/password"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	"github.com/kadekchresna/payroll/internal/v1/auth/dto"
	"github.com/kadekchresna/payroll/internal/v1/auth/model"
	auth_repository_interface "github.com/kadekchresna/payroll/internal/v1/auth/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/auth/usecase/interface"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface"
)

type userUsecase struct {
	userRepo     auth_repository_interface.IUserRepository
	employeeRepo repository_interface.IEmployeeRepository
	time         helper_time.TimeHelper
	config       config.Config
}

func NewUserUsecase(
	config config.Config,
	time helper_time.TimeHelper,
	userRepo auth_repository_interface.IUserRepository,
	employeeRepo repository_interface.IEmployeeRepository,

) usecase_interface.IUserUsecase {
	return &userUsecase{
		userRepo:     userRepo,
		time:         time,
		config:       config,
		employeeRepo: employeeRepo,
	}
}

func (uc *userUsecase) Create(ctx context.Context, req dto.CreateUserRequest) error {

	salt, _ := password.GenerateRandomSalt(16)
	hashed := password.HashPasswordWithSalt(req.Password, salt)
	u := &model.User{
		Username:  req.Username,
		Password:  hashed,
		Role:      "users",
		Status:    "active",
		Salt:      salt,
		CreatedBy: 0,
		UpdatedBy: 0,
		CreatedAt: uc.time.Now(),
		UpdatedAt: uc.time.Now(),
	}
	return uc.userRepo.Create(ctx, u)
}

func (uc *userUsecase) Login(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginUserResponse, error) {

	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) == 0 {
		return nil, errors.New("Username is required")
	}

	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, errors.New("Password is required")
	}

	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve user data, %s", err.Error())
	}

	if !password.ComparePasswordWithHash(req.Password, user.Salt, user.Password) {
		return nil, errors.New("Password is invalid")
	}

	employee, err := uc.employeeRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve employee data, %s", err.Error())
	}

	employeeID := 0
	employeeFullname := ""
	if employee != nil {
		employeeID = employee.ID
		employeeFullname = employee.FullName
	}

	accessToken, refreshToken, err := jwt.GenerateToken(uc.config.AppJWTSecret, user.ID, user.Role, employeeID, employeeFullname)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate token user, %s", err.Error())
	}

	return &dto.LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
