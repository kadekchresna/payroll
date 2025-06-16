package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	attendace_model "github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendace_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
	compensation_model "github.com/kadekchresna/payroll/internal/v1/compensation/model"
	overtimes_repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface"
	employee_repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface"
	"github.com/kadekchresna/payroll/internal/v1/payroll/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/payroll/usecase/interface"
	payslip_model "github.com/kadekchresna/payroll/internal/v1/payslip/model"
	payslip_repository_interface "github.com/kadekchresna/payroll/internal/v1/payslip/repository/interface"
)

type payrollUsecase struct {
	times                helper_time.TimeHelper
	overtimeRepo         overtimes_repository_interface.IOvertimeRepository
	reimbursementRepo    overtimes_repository_interface.IReimbursementRepository
	attendanceRepo       attendace_repository_interface.IAttendanceRepository
	attendancePeriodRepo attendace_repository_interface.IAttendancePeriodRepository
	auditRepo            audit_repository_interface.IAuditRepository
	employeeRepo         employee_repository_interface.IEmployeeRepository
	payslipRepo          payslip_repository_interface.IPayslipRepository
	transactionBundler   helper_db.ITransactionBundler
}

func NewPayrollUsecase(
	times helper_time.TimeHelper,
	overtimeRepo overtimes_repository_interface.IOvertimeRepository,
	reimbursementRepo overtimes_repository_interface.IReimbursementRepository,
	attendanceRepo attendace_repository_interface.IAttendanceRepository,
	attendancePeriodRepo attendace_repository_interface.IAttendancePeriodRepository,
	auditRepo audit_repository_interface.IAuditRepository,
	employeeRepo employee_repository_interface.IEmployeeRepository,
	payslipRepo payslip_repository_interface.IPayslipRepository,
	transactionBundler helper_db.ITransactionBundler,
) usecase_interface.IPayrollUsecase {
	return &payrollUsecase{
		times:                times,
		overtimeRepo:         overtimeRepo,
		reimbursementRepo:    reimbursementRepo,
		attendanceRepo:       attendanceRepo,
		auditRepo:            auditRepo,
		transactionBundler:   transactionBundler,
		attendancePeriodRepo: attendancePeriodRepo,
		employeeRepo:         employeeRepo,
		payslipRepo:          payslipRepo,
	}
}

func (r *payrollUsecase) CreatePayroll(ctx context.Context, req dto.CreatePayrollRequest) error {

	now := r.times.Now()

	if req.AttendancePeriodID <= 0 {
		return errors.New("Attendance period is required")
	}

	attendancePeriod, err := r.attendancePeriodRepo.GetByID(ctx, req.AttendancePeriodID)
	if err != nil {
		e := fmt.Errorf("Failed to retireve attendace period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendancePeriodRepo.GetByID. %s", e.Error()))
		return e
	}

	if attendancePeriod == nil {
		return errors.New("Attendance period is not found")
	}

	if attendancePeriod.IsPayslipGenerated {
		return errors.New("Payroll already executed for this attendance period")

	}

	attendanceEmployees, err := r.attendanceRepo.GetEmployeeCountByDateRange(ctx, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd)
	if err != nil {
		e := fmt.Errorf("Failed to retireve attendace per period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendanceRepo.GetEmployeeCountByDateRange. %s", e.Error()))
		return e
	}

	overtimeEmployees, err := r.overtimeRepo.SumOvertimeByDateRange(ctx, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd)
	if err != nil {
		e := fmt.Errorf("Failed to retireve overtime per period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error overtimeRepo.SumOvertimeByDateRange. %s", e.Error()))
		return e
	}

	reimbursements, err := r.reimbursementRepo.GetByDateRange(ctx, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd)
	if err != nil {
		e := fmt.Errorf("Failed to retireve reimbursement per period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error reimbursementRepo.GetByDateRange. %s", e.Error()))
		return e
	}

	reimbursementMappedEmployeeID := make(map[int][]compensation_model.Reimbursement, len(reimbursements))
	reimbursementIDs := make([]int, 0, len(reimbursements))

	for _, r := range reimbursements {
		reimbursementIDs = append(reimbursementIDs, r.ID)
		_, ok := reimbursementMappedEmployeeID[r.EmployeeID]
		if !ok {
			reimbursementMappedEmployeeID[r.EmployeeID] = []compensation_model.Reimbursement{r}
		} else {

			reimbursementMappedEmployeeID[r.EmployeeID] = append(reimbursementMappedEmployeeID[r.EmployeeID], r)
		}

	}

	reimbursementEmployees, err := r.reimbursementRepo.SumReimbursementsByID(ctx, reimbursementIDs)
	if err != nil {
		e := fmt.Errorf("Failed to retireve summary reimbursement per period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error reimbursementRepo.SumReimbursementsByID. %s", e.Error()))
		return e
	}

	payslipsMappedEmployeeID := map[int]payslip_model.Payslip{}

	for _, a := range attendanceEmployees {
		_, ok := payslipsMappedEmployeeID[a.EmployeeID]
		if !ok {
			payslipsMappedEmployeeID[a.EmployeeID] = payslip_model.Payslip{
				TotalAttendanceDays: a.Count,
			}
		} else {
			ps := payslipsMappedEmployeeID[a.EmployeeID]
			ps.TotalAttendanceDays += a.Count
			payslipsMappedEmployeeID[a.EmployeeID] = ps
		}
	}

	for _, o := range overtimeEmployees {
		_, ok := payslipsMappedEmployeeID[o.EmployeeID]
		if !ok {
			payslipsMappedEmployeeID[o.EmployeeID] = payslip_model.Payslip{
				TotalOvertimeHours: o.TotalHours,
			}
		} else {
			ps := payslipsMappedEmployeeID[o.EmployeeID]
			ps.TotalOvertimeHours += o.TotalHours
			payslipsMappedEmployeeID[o.EmployeeID] = ps
		}
	}

	for _, r := range reimbursementEmployees {
		_, ok := payslipsMappedEmployeeID[r.EmployeeID]
		if !ok {
			payslipsMappedEmployeeID[r.EmployeeID] = payslip_model.Payslip{
				TotalReimbursement: r.TotalAmount,
			}
		} else {
			ps := payslipsMappedEmployeeID[r.EmployeeID]
			ps.TotalReimbursement += r.TotalAmount
			payslipsMappedEmployeeID[r.EmployeeID] = ps
		}
	}

	workingDays := 0
	for d := attendancePeriod.PeriodStart; !d.After(attendancePeriod.PeriodEnd); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workingDays++
		}
	}

	err = r.transactionBundler.WithTransaction(ctx, func(ctx context.Context) error {

		for k := range payslipsMappedEmployeeID {
			employee, err := r.employeeRepo.GetByID(ctx, k)
			if err != nil {
				e := fmt.Errorf("Failed to retireve employee data, %s", err.Error())
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error employeeRepo.GetByID. %s", e.Error()))
				return e
			}

			p := payslipsMappedEmployeeID[k]

			p.CreatedAt = now
			p.UpdatedAt = now
			p.CreatedBy = req.UserID
			p.UpdatedBy = req.UserID
			p.PeriodID = attendancePeriod.ID

			salaryPerDay := employee.Salary / float64(workingDays)
			totalSalaryByAttendance := float64(p.TotalAttendanceDays) * math.Ceil(salaryPerDay)
			totalSalaryByOvertime := float64(p.TotalOvertimeHours) * (math.Ceil(salaryPerDay) / 8) * 2

			p.TotalAttendanceSalary = totalSalaryByAttendance
			p.TotalOvertimeSalary = totalSalaryByOvertime

			p.TotalTakeHomePay = p.TotalAttendanceSalary + p.TotalOvertimeSalary + p.TotalReimbursement
			p.EmployeeID = employee.ID

			payslipID, err := r.payslipRepo.Create(ctx, &p)
			if err != nil {
				e := fmt.Errorf("Failed to create payslip data, %s", err.Error())
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error employeeRepo.GetByID. %s", e.Error()))
				return e
			}

			rIds := make([]int, 0, len(reimbursementMappedEmployeeID[employee.ID]))
			for _, reimbursement := range reimbursementMappedEmployeeID[employee.ID] {
				rIds = append(rIds, reimbursement.ID)
				if err := r.auditRepo.Create(ctx, audit_model.AuditLog{
					TableName: "reimbursements",
					Action:    "update",
					RecordID:  reimbursement.ID,
					OldData: map[string]interface{}{
						"id":          reimbursement.ID,
						"employee_id": reimbursement.EmployeeID,
						"date":        reimbursement.Date,
						"amount":      reimbursement.Amount,
						"description": reimbursement.Description,
						"created_at":  reimbursement.UpdatedAt,
						"created_by":  reimbursement.UpdatedBy,
						"updated_at":  now,
						"updated_by":  req.UserID,
						"payslip_id":  payslipID,
					},
					NewData: map[string]interface{}{
						"id":          reimbursement.ID,
						"employee_id": reimbursement.EmployeeID,
						"date":        reimbursement.Date,
						"amount":      reimbursement.Amount,
						"description": reimbursement.Description,
						"created_at":  reimbursement.UpdatedAt,
						"updated_at":  reimbursement.CreatedBy,
						"created_by":  reimbursement.UpdatedBy,
						"updated_by":  reimbursement.UpdatedBy,
						"payslip_id":  reimbursement.PayslipID,
					},
					ChangedBy: req.UserID,
					ChangedAt: now,
				}); err != nil {
					e := fmt.Errorf("Failed to create reimbursement audit data, %s", err.Error())
					logger.LogWithContext(ctx).Error(fmt.Sprintf("error auditRepo.Create. %s", e.Error()))
					return e
				}

				if err := r.reimbursementRepo.Update(ctx, &compensation_model.Reimbursement{
					UpdatedAt: now,
					UpdatedBy: req.UserID,
					PayslipID: payslipID,
				}, rIds); err != nil {
					e := fmt.Errorf("Failed to update reimburse payslip data, %s", err.Error())
					logger.LogWithContext(ctx).Error(fmt.Sprintf("error reimbursementRepo.Update. %s", e.Error()))
					return e
				}
			}

			if err := r.auditRepo.Create(ctx, audit_model.AuditLog{
				TableName: "payslips",
				Action:    "create",
				RecordID:  payslipID,
				OldData:   map[string]interface{}{},
				NewData: map[string]interface{}{
					"id":                      payslipID,
					"employee_id":             p.EmployeeID,
					"total_attendance_days":   p.TotalAttendanceDays,
					"total_overtime_hours":    p.TotalOvertimeHours,
					"total_attendance_salary": p.TotalAttendanceSalary,
					"total_overtime_salary":   p.TotalOvertimeSalary,
					"total_reimbursement":     p.TotalReimbursement,
					"total_take_home_pay":     p.TotalTakeHomePay,
					"created_at":              p.CreatedAt,
					"updated_at":              p.UpdatedAt,
					"created_by":              p.CreatedBy,
					"updated_by":              p.UpdatedBy,
					"period_id":               p.PeriodID,
				},
				ChangedBy: p.CreatedBy,
				ChangedAt: now,
			}); err != nil {
				e := fmt.Errorf("Failed to create payslip audit data, %s", err.Error())
				logger.LogWithContext(ctx).Error(fmt.Sprintf("error auditRepo.Create. %s", e.Error()))
				return e
			}
		}

		ap := &attendace_model.AttendancePeriod{
			ID:                 attendancePeriod.ID,
			IsPayslipGenerated: true,
			UpdatedAt:          now,
			UpdatedBy:          req.UserID,
		}
		if err := r.attendancePeriodRepo.UpdatePeriod(ctx, ap); err != nil {
			e := fmt.Errorf("Failed to update attendance period data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendancePeriodRepo.UpdatePeriod. %s", e.Error()))
			return e
		}

		if err := r.auditRepo.Create(ctx, audit_model.AuditLog{
			TableName: "attendances_period",
			Action:    "update",
			RecordID:  attendancePeriod.ID,
			OldData: map[string]interface{}{
				"id":                   attendancePeriod.ID,
				"period_start":         attendancePeriod.PeriodStart,
				"period_end":           attendancePeriod.PeriodEnd,
				"is_payslip_generated": attendancePeriod.IsPayslipGenerated,
				"created_at":           attendancePeriod.CreatedAt,
				"updated_at":           attendancePeriod.UpdatedAt,
				"created_by":           attendancePeriod.CreatedBy,
				"updated_by":           attendancePeriod.UpdatedBy,
			},
			NewData: map[string]interface{}{
				"id":                   attendancePeriod.ID,
				"period_start":         attendancePeriod.PeriodStart,
				"period_end":           attendancePeriod.PeriodEnd,
				"is_payslip_generated": ap.IsPayslipGenerated,
				"created_at":           attendancePeriod.CreatedAt,
				"updated_at":           ap.UpdatedAt,
				"created_by":           attendancePeriod.CreatedBy,
				"updated_by":           ap.UpdatedBy,
			},
			ChangedBy: req.UserID,
			ChangedAt: now,
		}); err != nil {
			e := fmt.Errorf("Failed to update attendance audit data, %s", err.Error())
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

func (r *payrollUsecase) GetPayrollByID(ctx context.Context, req *dto.GetEmployeePayrollRequest) (*dto.GetEmployeePayrollResponse, error) {
	payslip, err := r.payslipRepo.GetByID(ctx, req.PayslipID)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve payslip data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error payslipRepo.GetByID. %s", e.Error()))
		return nil, e
	}

	if payslip.EmployeeID != req.EmployeeID {
		return nil, errors.New("Payslip is not found")
	}

	attendancePeriod, err := r.attendancePeriodRepo.GetByID(ctx, payslip.PeriodID)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve attendance period data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error attendancePeriodRepo.GetByID. %s", e.Error()))
		return nil, e
	}

	em, err := r.employeeRepo.GetByID(ctx, payslip.EmployeeID)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve employee data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error employeeRepo.GetByID. %s", e.Error()))
		return nil, e
	}

	reimbursements, err := r.reimbursementRepo.GetByPayslipID(ctx, payslip.ID, payslip.EmployeeID)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve reimbursements data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error reimbursementRepo.GetByPayslipID. %s", e.Error()))
		return nil, e
	}

	reimbursementsRes := make([]dto.Reimbursement, 0, len(reimbursements))
	for _, da := range reimbursements {
		reimbursementsRes = append(reimbursementsRes, dto.Reimbursement{
			Date:   da.Date,
			Amount: da.Amount,
		})
	}

	workingDays := 0
	for d := attendancePeriod.PeriodStart; !d.After(attendancePeriod.PeriodEnd); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workingDays++
		}
	}

	return &dto.GetEmployeePayrollResponse{
		TotalWorkingDaysInAttendancePeriod: workingDays,
		TotalSalaryPerMonth:                em.Salary,
		TotalAttendanceDays:                payslip.TotalAttendanceDays,
		TotalOvertimeHours:                 payslip.TotalOvertimeHours,
		TotalAttendanceSalary:              payslip.TotalAttendanceSalary,
		TotalOvertimeSalary:                payslip.TotalOvertimeSalary,
		TotalReimbursement:                 payslip.TotalReimbursement,
		TotalTakeHomePay:                   payslip.TotalTakeHomePay,
		PeriodStart:                        attendancePeriod.PeriodStart,
		PeriodEnd:                          attendancePeriod.PeriodEnd,
		Reimbursement:                      reimbursementsRes,
	}, nil
}

func (r *payrollUsecase) GetPayrollSummary(ctx context.Context) (*dto.GetPayrollSummary, error) {

	thpEmplyees, err := r.payslipRepo.GetTotalTakeHomePayPerEmployee(ctx)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve thp employee data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error payslipRepo.GetTotalTakeHomePayPerEmployee. %s", e.Error()))
		return nil, e
	}

	allThp, err := r.payslipRepo.GetTotalTakeHomePayAllEmployees(ctx)
	if err != nil {
		e := fmt.Errorf("Failed to retrieve all thp data, %s", err.Error())
		logger.LogWithContext(ctx).Error(fmt.Sprintf("error payslipRepo.GetTotalTakeHomePayAllEmployees. %s", e.Error()))
		return nil, e
	}

	res := make([]dto.EmployeePayroll, 0, len(thpEmplyees))
	for _, e := range thpEmplyees {
		res = append(res, dto.EmployeePayroll{
			EmployeeID:       e.EmployeeID,
			TotalTakeHomePay: e.TotalTakeHomePay,
		})
	}

	return &dto.GetPayrollSummary{
		EmployeePayrollSummary: res,
		TotalTakeHomePayAll:    allThp,
	}, nil
}
