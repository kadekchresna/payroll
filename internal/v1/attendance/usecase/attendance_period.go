package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendance_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
)

type attendancePeriodUsecase struct {
	times                helper_time.TimeHelper
	attendancePeriodRepo attendance_repository_interface.IAttendancePeriodRepository
	transactionBundler   helper_db.ITransactionBundler
	auditRepo            audit_repository_interface.IAuditRepository
}

func NewAttendancePeriodUsecase(
	times helper_time.TimeHelper,
	attendancePeriodRepo attendance_repository_interface.IAttendancePeriodRepository,
	transactionBundler helper_db.ITransactionBundler,
	auditRepo audit_repository_interface.IAuditRepository,
) usecase_interface.IAttendancePeriodUsecase {
	return &attendancePeriodUsecase{
		times:                times,
		attendancePeriodRepo: attendancePeriodRepo,
		transactionBundler:   transactionBundler,
		auditRepo:            auditRepo,
	}
}

func (u *attendancePeriodUsecase) CreateAttendancePeriod(ctx context.Context, req *dto.CreateAttendancePeriodRequest) error {

	now := u.times.Now()

	if req.PeriodStart.IsZero() {
		return errors.New("Period start is required")
	}

	if req.PeriodEnd.IsZero() {
		return errors.New("Period end is required")
	}

	if req.PeriodStart.After(req.PeriodEnd) {
		return errors.New("Period start must be behind period end")
	}

	existingPeriod, err := u.attendancePeriodRepo.GetByPeriodIntersect(ctx, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		e := fmt.Errorf("Failed to retieve attendance period data. %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendancePeriodRepo.GetByPeriodIntersect. %s", e.Error()))
		return e
	}

	if existingPeriod != nil {
		return errors.New("Attendance period is already exist for specified period range")
	}

	err = u.transactionBundler.WithTransaction(ctx, func(ctx context.Context) error {

		ap := &model.AttendancePeriod{
			PeriodStart:        req.PeriodStart,
			PeriodEnd:          req.PeriodEnd,
			IsPayslipGenerated: false,
			CreatedAt:          now,
			UpdatedAt:          now,
			UpdatedBy:          req.UserID,
			CreatedBy:          req.UserID,
		}

		attendancePeriodID, err := u.attendancePeriodRepo.Create(ctx, ap)
		if err != nil {
			e := fmt.Errorf("Failed to retieve attendance period data. %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendancePeriodRepo.GetByPeriodIntersect. %s", e.Error()))
			return e
		}

		if err := u.auditRepo.Create(ctx, audit_model.AuditLog{
			TableName: "attendances_period",
			Action:    "create",
			RecordID:  attendancePeriodID,
			OldData:   map[string]interface{}{},
			NewData: map[string]interface{}{
				"id":                   attendancePeriodID,
				"period_start":         ap.PeriodStart,
				"period_end":           ap.PeriodEnd,
				"is_payslip_generated": ap.IsPayslipGenerated,
				"created_at":           ap.CreatedAt,
				"updated_at":           ap.UpdatedAt,
				"created_by":           ap.CreatedBy,
				"updated_by":           ap.UpdatedBy,
			},
			ChangedBy: ap.CreatedBy,
			ChangedAt: now,
		}); err != nil {
			e := fmt.Errorf("Failed to create attendance audit data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error auditRepo.Create. %s", e.Error()))
			return e
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
