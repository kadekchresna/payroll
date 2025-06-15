package repository

import (
	"context"
	"time"

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

// --- helper conversion ---
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

func toDAO(e *model.Employee) *dao.Employee {
	return &dao.Employee{
		ID:        e.ID,
		FullName:  e.FullName,
		Salary:    e.Salary,
		Code:      e.Code,
		UserID:    e.UserID,
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
}

// --- CRUD implementations ---

func (r *employeeRepo) Create(ctx context.Context, e *model.Employee) error {
	daoEmp := toDAO(e)
	return r.db.WithContext(ctx).Create(daoEmp).Error
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

func (r *employeeRepo) Update(ctx context.Context, e *model.Employee) error {
	daoEmp := toDAO(e)
	daoEmp.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(daoEmp).Error
}

func (r *employeeRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&dao.Employee{}, id).Error
}

func (r *employeeRepo) ListAll(ctx context.Context) ([]model.Employee, error) {
	var daos []dao.Employee
	if err := r.db.WithContext(ctx).Find(&daos).Error; err != nil {
		return nil, err
	}

	var result []model.Employee
	for _, d := range daos {
		result = append(result, *toModel(&d))
	}
	return result, nil
}
