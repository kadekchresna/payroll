package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
	attendance_model "github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendance_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
	employee_repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface"
)

type attendanceUsecase struct {
	attendanceRepo     attendance_repository_interface.IAttendanceRepository
	employeeRepo       employee_repository_interface.IEmployeeRepository
	auditRepo          audit_repository_interface.IAuditRepository
	transactionBundler helper_db.ITransactionBundler
	times              helper_time.TimeHelper
}

func NewAttendanceUsecase(
	times helper_time.TimeHelper,
	attendanceRepo attendance_repository_interface.IAttendanceRepository,
	employeeRepo employee_repository_interface.IEmployeeRepository,
	auditRepo audit_repository_interface.IAuditRepository,
	transactionBundler helper_db.ITransactionBundler,
) usecase_interface.IAttendanceUsecase {
	return &attendanceUsecase{
		times:              times,
		attendanceRepo:     attendanceRepo,
		employeeRepo:       employeeRepo,
		auditRepo:          auditRepo,
		transactionBundler: transactionBundler,
	}
}

func (u *attendanceUsecase) CreateAttendance(ctx context.Context, req *dto.CreateAttendanceRequest) error {

	now := u.times.Now()

	if helper_time.IsWeekend(now.Weekday()) {
		return errors.New("Today is weekend")
	}

	if req.EmployeeID == 0 {
		return errors.New("Employee ID is empty")
	}

	// check employee id exist or not
	employee, err := u.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		e := fmt.Errorf("Failed to retireve employee data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error employeeRepo.GetByID. %s", e.Error()))
		return e
	}

	isCheckOut := false
	attendace, err := u.attendanceRepo.GetByDateAndEmployeeID(ctx, employee.ID, req.Date)
	if err != nil {
		e := fmt.Errorf("Failed to retireve attendace data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendanceRepo.GetByDateAndEmployeeID. %s", e.Error()))
		return e
	}

	if attendace != nil {
		isCheckOut = true
	}

	if attendace != nil && attendace.CheckedOutAt != nil {
		return nil
	}

	if err = u.transactionBundler.WithTransaction(ctx, func(ctx context.Context) error {

		a := &attendance_model.Attendance{
			EmployeeID:  employee.ID,
			Date:        req.Date,
			CreatedBy:   req.UserID,
			UpdatedBy:   req.UserID,
			CreatedAt:   now,
			UpdatedAt:   now,
			CheckedInAt: now,
		}

		// insert attendance along with audit log
		attendanceID, err := u.attendanceRepo.Create(ctx, a)
		if err != nil {
			e := fmt.Errorf("Failed to create attendance data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendanceRepo.Create. %s", e.Error()))
			return e
		}

		actionAudit := "create"
		oldData := map[string]interface{}{}
		newData := map[string]interface{}{
			"id":             attendanceID,
			"employee_id":    a.EmployeeID,
			"date":           a.Date,
			"created_at":     a.CreatedAt,
			"created_by":     a.CreatedBy,
			"updated_at":     a.UpdatedAt,
			"updated_by":     a.UpdatedBy,
			"checked_in_at":  a.CheckedInAt,
			"checked_out_at": a.CheckedOutAt,
		}
		if isCheckOut {
			actionAudit = "update"
			oldData = map[string]interface{}{
				"id":             attendanceID,
				"employee_id":    attendace.EmployeeID,
				"date":           attendace.Date,
				"created_at":     attendace.CreatedAt,
				"created_by":     attendace.CreatedBy,
				"updated_at":     attendace.UpdatedAt,
				"updated_by":     attendace.UpdatedBy,
				"checked_in_at":  attendace.CheckedInAt,
				"checked_out_at": attendace.CheckedOutAt,
			}

			newData = map[string]interface{}{
				"id":             attendanceID,
				"employee_id":    attendace.EmployeeID,
				"date":           attendace.Date,
				"created_at":     attendace.CreatedAt,
				"created_by":     attendace.CreatedBy,
				"updated_at":     a.UpdatedAt,
				"updated_by":     a.UpdatedBy,
				"checked_in_at":  attendace.CheckedInAt,
				"checked_out_at": now,
			}

		}

		// insert attendance along with audit log
		if err := u.auditRepo.Create(ctx, audit_model.AuditLog{
			TableName: "attendances",
			Action:    actionAudit,
			RecordID:  attendanceID,
			OldData:   oldData,
			NewData:   newData,
			ChangedAt: now,
			ChangedBy: req.UserID,
		}); err != nil {
			e := fmt.Errorf("Failed to create attendance audit data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error auditRepo.Create. %s", e.Error()))
			return e
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
