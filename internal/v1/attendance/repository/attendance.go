package repository

import (
	"context"

	"errors"
	"time"

	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	"github.com/kadekchresna/payroll/internal/v1/attendance/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) repository_interface.IAttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

func (r *attendanceRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := helper_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *attendanceRepository) Create(ctx context.Context, a *model.Attendance) (int, error) {

	db := r.getDB(ctx)
	da := dao.AttendanceDAO{
		EmployeeID:  a.EmployeeID,
		Date:        a.Date,
		CreatedBy:   a.CreatedBy,
		UpdatedBy:   a.UpdatedBy,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		CheckedInAt: a.CheckedInAt,
	}

	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "employee_id"}, {Name: "date"}}, // define conflict columns
		DoUpdates: clause.Assignments(map[string]interface{}{
			"updated_by":     a.UpdatedBy,
			"updated_at":     a.UpdatedAt,
			"checked_out_at": a.UpdatedAt,
		}),
	}).Create(&da).Error

	if err != nil {
		return 0, err
	}

	return da.ID, nil
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

	checkedOutAtValue, _ := da.CheckedOutAt.Value()
	var checkedOutAt *time.Time
	if t, ok := checkedOutAtValue.(time.Time); ok {
		checkedOutAt = &t
	} else if t, ok := checkedOutAtValue.(*time.Time); ok {
		checkedOutAt = t
	}
	return &model.Attendance{
		ID:           da.ID,
		EmployeeID:   da.EmployeeID,
		Date:         da.Date,
		CreatedAt:    da.CreatedAt,
		UpdatedAt:    da.UpdatedAt,
		CreatedBy:    da.CreatedBy,
		UpdatedBy:    da.UpdatedBy,
		CheckedInAt:  da.CheckedInAt,
		CheckedOutAt: checkedOutAt,
	}, nil
}

func (r *attendanceRepository) GetByDateAndEmployeeID(ctx context.Context, employeeID int, date time.Time) (*model.Attendance, error) {
	var da dao.AttendanceDAO

	db := r.db.WithContext(ctx)

	if employeeID > 0 {
		db = db.Where("employee_id = ?", employeeID)
	}

	if !date.IsZero() {
		db = db.Where("date = ?", date.Format(time.DateOnly))
	}

	err := db.First(&da).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	checkedOutAtValue, _ := da.CheckedOutAt.Value()
	var checkedOutAt *time.Time
	if t, ok := checkedOutAtValue.(time.Time); ok {
		checkedOutAt = &t
	} else if t, ok := checkedOutAtValue.(*time.Time); ok {
		checkedOutAt = t
	}
	return &model.Attendance{
		ID:           da.ID,
		EmployeeID:   da.EmployeeID,
		Date:         da.Date,
		CreatedAt:    da.CreatedAt,
		UpdatedAt:    da.UpdatedAt,
		CreatedBy:    da.CreatedBy,
		UpdatedBy:    da.UpdatedBy,
		CheckedInAt:  da.CheckedInAt,
		CheckedOutAt: checkedOutAt,
	}, nil
}

func (r *attendanceRepository) GetEmployeeCountByDateRange(ctx context.Context, periodStart time.Time, periodEnd time.Time) ([]*model.EmployeeAttendanceCount, error) {
	db := r.getDB(ctx)

	var daos []dao.AttendanceCountDAO
	query := `
		WITH attandances_period AS (
			SELECT employee_id 
			FROM attendances 
			WHERE date >= ? AND date <= ?
		)
		SELECT employee_id, COUNT(employee_id) 
		FROM attandances_period 
		GROUP BY employee_id;
	`

	if err := db.Raw(query, periodStart, periodEnd).Scan(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]*model.EmployeeAttendanceCount, 0, len(daos))
	for _, daoItem := range daos {
		result = append(result, &model.EmployeeAttendanceCount{
			EmployeeID: daoItem.EmployeeID,
			Count:      daoItem.Count,
		})
	}

	return result, nil
}
