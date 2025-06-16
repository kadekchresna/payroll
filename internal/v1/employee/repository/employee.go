package repository

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/employee/model"
	"github.com/kadekchresna/payroll/internal/v1/employee/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface"
	"gorm.io/gorm"
)

type employeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) repository_interface.IEmployeeRepository {
	return &employeeRepo{db: db}
}

func toModel(e *dao.Employee) *model.Employee {
	return &model.Employee{
		ID:        e.ID,
		FullName:  e.FullName,
		Salary:    e.Salary,
		Code:      e.Code,
		UserID:    e.UserID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
}

func (r *employeeRepo) GetByID(ctx context.Context, id int) (*model.Employee, error) {
	var emp dao.Employee
	if err := r.db.WithContext(ctx).First(&emp, id).Error; err != nil {
		return nil, err
	}
	return toModel(&emp), nil
}

func (r *employeeRepo) GetByUserID(ctx context.Context, userID int) (*model.Employee, error) {
	var e dao.Employee
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&e).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toModel(&e), nil
}
