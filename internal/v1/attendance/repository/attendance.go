package repository

import (
	"context"
	"errors"

	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	"github.com/kadekchresna/payroll/internal/v1/attendance/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	"gorm.io/gorm"
)

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) repository_interface.IAttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

func (r *attendanceRepository) Create(ctx context.Context, a *model.Attendance) error {
	da := dao.AttendanceDAO{
		EmployeeID: a.EmployeeID,
		Date:       a.Date,
		CreatedBy:  a.CreatedBy,
		UpdatedBy:  a.UpdatedBy,
	}
	return r.db.WithContext(ctx).Create(&da).Error
}

func (r *attendanceRepository) GetByID(ctx context.Context, id int) (*model.Attendance, error) {
	var da dao.AttendanceDAO
	err := r.db.WithContext(ctx).First(&da, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model.Attendance{
		ID:         da.ID,
		EmployeeID: da.EmployeeID,
		Date:       da.Date,
		CreatedAt:  da.CreatedAt,
		UpdatedAt:  da.UpdatedAt,
		CreatedBy:  da.CreatedBy,
		UpdatedBy:  da.UpdatedBy,
	}, nil
}

func (r *attendanceRepository) Update(ctx context.Context, a *model.Attendance) error {
	return r.db.WithContext(ctx).
		Model(&dao.AttendanceDAO{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"employee_id": a.EmployeeID,
			"date":        a.Date,
			"updated_by":  a.UpdatedBy,
		}).Error
}

func (r *attendanceRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&dao.AttendanceDAO{}, id).Error
}
