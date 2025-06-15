package repository

import (
	"context"
	"errors"
	"time"

	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
	"github.com/kadekchresna/payroll/internal/v1/compensation/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface"
	"gorm.io/gorm"
)

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) repository_interface.IOvertimeRepository {
	return &overtimeRepository{db: db}
}

func (r *overtimeRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := helper_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *overtimeRepository) Create(ctx context.Context, ot *model.Overtime) (int, error) {

	db := r.getDB(ctx)
	da := dao.OvertimeDAO{
		EmployeeID: ot.EmployeeID,
		Date:       ot.Date,
		Hours:      ot.Hours,
		CreatedBy:  ot.CreatedBy,
		UpdatedBy:  ot.UpdatedBy,
		CreatedAt:  ot.CreatedAt,
		UpdatedAt:  ot.UpdatedAt,
	}

	if err := db.WithContext(ctx).Create(&da).Error; err != nil {
		return 0, err
	}
	return da.ID, nil
}

func (r *overtimeRepository) GetByDateAndEmployeeID(ctx context.Context, employeeID int, date time.Time) ([]model.Overtime, error) {
	var os []dao.OvertimeDAO
	if err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Where("date = ?", date.Format(time.DateOnly)).Find(&os).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	res := make([]model.Overtime, len(os))

	for _, o := range os {
		res = append(res, model.Overtime{
			ID:         o.ID,
			EmployeeID: o.EmployeeID,
			Date:       o.Date,
			Hours:      o.Hours,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
			CreatedBy:  o.CreatedBy,
			UpdatedBy:  o.UpdatedBy,
		})
	}

	return res, nil
}
