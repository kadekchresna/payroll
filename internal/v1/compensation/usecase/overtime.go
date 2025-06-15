package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	attendace_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
	overtimes_repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/compensation/usecase/interface"
)

type overtimeUsecase struct {
	times              helper_time.TimeHelper
	overtimeRepo       overtimes_repository_interface.IOvertimeRepository
	attendanceRepo     attendace_repository_interface.IAttendanceRepository
	auditRepo          audit_repository_interface.IAuditRepository
	transactionBundler helper_db.ITransactionBundler
}

func NewOvertimeUsecase(
	overtimeRepo overtimes_repository_interface.IOvertimeRepository,
	attendanceRepo attendace_repository_interface.IAttendanceRepository,
	auditRepo audit_repository_interface.IAuditRepository,
	transactionBundler helper_db.ITransactionBundler,
	times helper_time.TimeHelper,
) usecase_interface.IOvertimeUsecase {
	return &overtimeUsecase{
		overtimeRepo:       overtimeRepo,
		attendanceRepo:     attendanceRepo,
		auditRepo:          auditRepo,
		transactionBundler: transactionBundler,
		times:              times,
	}
}

func (u *overtimeUsecase) CreateOvertime(ctx context.Context, req *dto.CreateOvertimeRequest) error {

	if req.EmployeeID == 0 {
		return errors.New("Employee ID is required")
	}

	if req.Hours <= 0 {
		return errors.New("Overtime hour must be inputed")
	}

	now := u.times.Now()

	latestAttendance, err := u.attendanceRepo.GetByDateAndEmployeeID(ctx, req.EmployeeID, req.Date)
	if err != nil {
		e := fmt.Errorf("Failed to retreive attendance data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendanceRepo.GetByDateAndEmployeeID. %s", e.Error()))
		return e
	}

	if latestAttendance.CheckedOutAt == nil {
		return fmt.Errorf("You are not yet clocked out at the specified date %s", req.Date.Format(time.DateOnly))
	}

	overtimes, err := u.overtimeRepo.GetByDateAndEmployeeID(ctx, req.EmployeeID, req.Date)
	if err != nil {
		e := fmt.Errorf("Failed to retreive overtime data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error overtimeRepo.GetByDateAndEmployeeID. %s", e.Error()))
		return e
	}

	totalHourOvertimePerDate := 0
	for _, o := range overtimes {
		totalHourOvertimePerDate += o.Hours
	}

	if totalHourOvertimePerDate+req.Hours > 3 {
		return errors.New("Overtime hours are above the limit per day")
	}

	err = u.transactionBundler.WithTransaction(ctx, func(ctx context.Context) error {
		ot := &model.Overtime{
			EmployeeID: req.EmployeeID,
			Date:       req.Date,
			Hours:      req.Hours,
			CreatedAt:  now,
			UpdatedAt:  now,
			CreatedBy:  req.UserID,
			UpdatedBy:  req.UserID,
		}
		overtimeID, err := u.overtimeRepo.Create(ctx, ot)
		if err != nil {
			e := fmt.Errorf("Failed to create overtime data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error overtimeRepo.Create. %s", e.Error()))
			return e
		}

		if err := u.auditRepo.Create(ctx, audit_model.AuditLog{
			TableName: "overtimes",
			Action:    "create",
			RecordID:  overtimeID,
			OldData:   map[string]interface{}{},
			NewData: map[string]interface{}{
				"id":          overtimeID,
				"employee_id": ot.EmployeeID,
				"date":        ot.Date,
				"hours":       ot.Hours,
				"created_at":  ot.CreatedAt,
				"updated_at":  ot.UpdatedAt,
				"created_by":  ot.CreatedBy,
				"updated_by":  ot.UpdatedBy,
			},
			ChangedBy: ot.CreatedBy,
			ChangedAt: now,
		}); err != nil {
			e := fmt.Errorf("Failed to create overtime audit data, %s", err.Error())
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
