package repository

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/auth/model"
	"github.com/kadekchresna/payroll/internal/v1/auth/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/auth/repository/interface"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository_interface.IUserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	ud := dao.UserDAO{
		Username:  user.Username,
		Password:  user.Password,
		Salt:      user.Salt,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	}
	return r.db.WithContext(ctx).Create(&ud).Error
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
