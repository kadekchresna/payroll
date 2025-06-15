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
)

type attendancePeriodRepository struct {
	db *gorm.DB
}

func NewAttendancePeriodRepository(db *gorm.DB) repository_interface.IAttendancePeriodRepository {
	return &attendancePeriodRepository{
		db: db,
	}
}

func (r *attendancePeriodRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := helper_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *attendancePeriodRepository) Create(ctx context.Context, p *model.AttendancePeriod) (int, error) {
	db := r.getDB(ctx)
	daoPeriod := dao.AttendancePeriodDAO{
		PeriodStart:        p.PeriodStart,
		PeriodEnd:          p.PeriodEnd,
		IsPayslipGenerated: &p.IsPayslipGenerated,
		CreatedBy:          p.CreatedBy,
		UpdatedBy:          p.UpdatedBy,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}

	if err := db.WithContext(ctx).Create(&daoPeriod).Error; err != nil {
		return 0, err
	}

	return daoPeriod.ID, nil
}

func (r *attendancePeriodRepository) GetByPeriodIntersect(ctx context.Context, periodStart time.Time, periodEnd time.Time) (*model.AttendancePeriod, error) {
	var daoPeriod dao.AttendancePeriodDAO

	if err := r.db.WithContext(ctx).
		Where("period_start <= ? AND period_end >= ?", periodEnd, periodEnd).
		Or("period_start <= ? AND period_end >= ?", periodStart, periodStart).
		Or("period_start <= ? AND period_end >= ?", periodStart, periodEnd).
		Or("period_start >= ? AND period_end <= ?", periodStart, periodEnd).
		First(&daoPeriod).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &model.AttendancePeriod{
		ID:                 daoPeriod.ID,
		PeriodStart:        daoPeriod.PeriodStart,
		PeriodEnd:          daoPeriod.PeriodEnd,
		IsPayslipGenerated: getBoolValue(daoPeriod.IsPayslipGenerated),
	}, nil
}

func getBoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func (r *attendancePeriodRepository) GetByID(ctx context.Context, id int) (*model.AttendancePeriod, error) {
	var daoPeriod dao.AttendancePeriodDAO

	if err := r.db.WithContext(ctx).
		First(&daoPeriod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &model.AttendancePeriod{
		ID:                 daoPeriod.ID,
		PeriodStart:        daoPeriod.PeriodStart,
		PeriodEnd:          daoPeriod.PeriodEnd,
		IsPayslipGenerated: getBoolValue(daoPeriod.IsPayslipGenerated),
	}, nil
}

func (r *attendancePeriodRepository) UpdatePeriod(ctx context.Context, ap *model.AttendancePeriod) error {
	db := r.getDB(ctx)

	updateData := map[string]interface{}{
		"is_payslip_generated": ap.IsPayslipGenerated,
		"updated_at":           ap.UpdatedAt,
		"updated_by":           ap.UpdatedBy,
	}

	return db.WithContext(ctx).Model(&dao.AttendancePeriodDAO{}).
		Where("id = ?", ap.ID).
		Updates(updateData).Error
}
