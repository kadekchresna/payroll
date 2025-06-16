package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db_test "github.com/kadekchresna/payroll/infrastructure/db/helper/mocks"
	attendace_model "github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendace_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface/mocks"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface/mocks"
	compensation_model "github.com/kadekchresna/payroll/internal/v1/compensation/model"
	compensation_repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface/mocks"
	employee_model "github.com/kadekchresna/payroll/internal/v1/employee/model"
	employee_repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface/mocks"
	"github.com/kadekchresna/payroll/internal/v1/payroll/dto"
	payslip_model "github.com/kadekchresna/payroll/internal/v1/payslip/model"
	payslip_repository_interface "github.com/kadekchresna/payroll/internal/v1/payslip/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_payrollUsecase_GetPayrollSummary(t *testing.T) {

	res := dto.GetPayrollSummary{
		EmployeePayrollSummary: []dto.EmployeePayroll{
			{
				EmployeeID:       1,
				TotalTakeHomePay: 10000,
			},
		},
		TotalTakeHomePayAll: 10000,
	}

	now := time.Now()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		want       *dto.GetPayrollSummary
		wantErr    bool
		beforeFunc func(
			attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
			attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
			overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
			reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			employeeRepo *employee_repository_interface.MockIEmployeeRepository,
			payslipRepo *payslip_repository_interface.MockIPayslipRepository,
			transactionBundler *helper_db_test.MockITransactionBundler,
		) helper_time.TimeHelper
	}{
		{
			name: "GetPayrollSummary-Success",
			args: args{
				ctx: context.Background(),
			},
			want:    &res,
			wantErr: false,
			beforeFunc: func(
				attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
				attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				payslipRepo *payslip_repository_interface.MockIPayslipRepository,
				transactionBundler *helper_db_test.MockITransactionBundler,
			) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				payslipRepo.EXPECT().GetTotalTakeHomePayPerEmployee(mock.Anything).Return([]payslip_model.TotalTakeHomePayPerEmployee{
					{
						EmployeeID:       1,
						TotalTakeHomePay: 10000,
					},
				}, nil).Once()

				payslipRepo.EXPECT().GetTotalTakeHomePayAllEmployees(mock.Anything).Return(10000, nil).Once()

				return times
			},
		},
		{
			name: "GetPayrollSummary-FailedErrorQueryGetTotalTakeHomePayAllEmployees",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(
				attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
				attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				payslipRepo *payslip_repository_interface.MockIPayslipRepository,
				transactionBundler *helper_db_test.MockITransactionBundler,
			) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				payslipRepo.EXPECT().GetTotalTakeHomePayPerEmployee(mock.Anything).Return([]payslip_model.TotalTakeHomePayPerEmployee{
					{
						EmployeeID:       1,
						TotalTakeHomePay: 10000,
					},
				}, nil).Once()

				payslipRepo.EXPECT().GetTotalTakeHomePayAllEmployees(mock.Anything).Return(0, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "GetPayrollSummary-FailedErrorQueryGetTotalTakeHomePayPerEmployee",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(
				attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
				attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				payslipRepo *payslip_repository_interface.MockIPayslipRepository,
				transactionBundler *helper_db_test.MockITransactionBundler,
			) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				payslipRepo.EXPECT().GetTotalTakeHomePayPerEmployee(mock.Anything).Return(nil, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendancePeriodRepo := attendace_repository_interface.NewMockIAttendancePeriodRepository(t)
			attendanceRepo := attendace_repository_interface.NewMockIAttendanceRepository(t)
			overtimeRepo := compensation_repository_interface.NewMockIOvertimeRepository(t)
			reimbursementRepo := compensation_repository_interface.NewMockIReimbursementRepository(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			employeeRepo := employee_repository_interface.NewMockIEmployeeRepository(t)
			payslipRepo := payslip_repository_interface.NewMockIPayslipRepository(t)
			transactionBundler := helper_db_test.NewMockITransactionBundler(t)

			times := tt.beforeFunc(
				attendancePeriodRepo,
				attendanceRepo,
				overtimeRepo,
				reimbursementRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)

			r := NewPayrollUsecase(
				times,
				overtimeRepo,
				reimbursementRepo,
				attendanceRepo,
				attendancePeriodRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)
			got, err := r.GetPayrollSummary(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("payrollUsecase.GetPayrollSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payrollUsecase.GetPayrollSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payrollUsecase_GetPayrollByID(t *testing.T) {

	now := time.Now()

	req := dto.GetEmployeePayrollRequest{
		EmployeeID: 1,
		PayslipID:  1,
	}

	res := dto.GetEmployeePayrollResponse{
		TotalWorkingDaysInAttendancePeriod: 2,
		TotalSalaryPerMonth:                10000,
		TotalAttendanceDays:                1,
		TotalOvertimeHours:                 1,
		TotalAttendanceSalary:              500,
		TotalOvertimeSalary:                62.5,
		TotalReimbursement:                 10000,
		TotalTakeHomePay:                   10562.5,
		PeriodStart:                        now.Add(-24 * time.Hour),
		PeriodEnd:                          now.Add(24 * time.Hour),
		Reimbursement: []dto.Reimbursement{
			{
				Date:   now,
				Amount: 10000,
			},
		},
	}

	p := payslip_model.Payslip{
		ID:                    1,
		EmployeeID:            1,
		TotalAttendanceDays:   res.TotalAttendanceDays,
		TotalOvertimeHours:    res.TotalOvertimeHours,
		TotalAttendanceSalary: res.TotalAttendanceSalary,
		TotalOvertimeSalary:   res.TotalOvertimeSalary,
		TotalReimbursement:    res.TotalReimbursement,
		TotalTakeHomePay:      res.TotalTakeHomePay,
		PeriodID:              1,
	}

	attendancePeriod := attendace_model.AttendancePeriod{
		ID:                 1,
		PeriodStart:        res.PeriodStart,
		PeriodEnd:          res.PeriodEnd,
		IsPayslipGenerated: true,
	}

	r := compensation_model.Reimbursement{
		ID:     1,
		Date:   now,
		Amount: 10000,
	}

	type args struct {
		ctx context.Context
		req *dto.GetEmployeePayrollRequest
	}
	tests := []struct {
		name       string
		args       args
		want       *dto.GetEmployeePayrollResponse
		wantErr    bool
		beforeFunc func(
			attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
			attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
			overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
			reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			employeeRepo *employee_repository_interface.MockIEmployeeRepository,
			payslipRepo *payslip_repository_interface.MockIPayslipRepository,
			transactionBundler *helper_db_test.MockITransactionBundler,
		) helper_time.TimeHelper
	}{
		{
			name: "GetPayrollByID-Success",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    &res,
			wantErr: false,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&p, nil).Once()

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()

				employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{Salary: 10000}, nil).Once()

				reimbursementRepo.EXPECT().GetByPayslipID(mock.Anything, p.ID, p.EmployeeID).Return([]compensation_model.Reimbursement{r}, nil).Once()

				return times

			},
		},

		{
			name: "GetPayrollByID-FailedGetByPayslipID",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&p, nil).Once()

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()

				employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{Salary: 10000}, nil).Once()

				reimbursementRepo.EXPECT().GetByPayslipID(mock.Anything, p.ID, p.EmployeeID).Return([]compensation_model.Reimbursement{r}, errors.New("FATAL ERROR")).Once()

				return times

			},
		},

		{
			name: "GetPayrollByID-FailedGetByID",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&p, nil).Once()

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()

				employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{Salary: 10000}, errors.New("FATAL ERROR")).Once()

				return times

			},
		},

		{
			name: "GetPayrollByID-FailedGetByID",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&p, nil).Once()

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, errors.New("FATAL ERROR")).Once()

				return times

			},
		},

		{
			name: "GetPayrollByID-FailedGetByID",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&p, errors.New("FATAL ERROR")).Once()

				return times

			},
		},

		{
			name: "GetPayrollByID-FailedEmployeeIDDifferentPayslip",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				ps := p
				ps.EmployeeID = 12
				payslipRepo.EXPECT().GetByID(mock.Anything, req.PayslipID).Return(&ps, nil).Once()

				return times

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendancePeriodRepo := attendace_repository_interface.NewMockIAttendancePeriodRepository(t)
			attendanceRepo := attendace_repository_interface.NewMockIAttendanceRepository(t)
			overtimeRepo := compensation_repository_interface.NewMockIOvertimeRepository(t)
			reimbursementRepo := compensation_repository_interface.NewMockIReimbursementRepository(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			employeeRepo := employee_repository_interface.NewMockIEmployeeRepository(t)
			payslipRepo := payslip_repository_interface.NewMockIPayslipRepository(t)
			transactionBundler := helper_db_test.NewMockITransactionBundler(t)

			times := tt.beforeFunc(
				attendancePeriodRepo,
				attendanceRepo,
				overtimeRepo,
				reimbursementRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)

			r := NewPayrollUsecase(
				times,
				overtimeRepo,
				reimbursementRepo,
				attendanceRepo,
				attendancePeriodRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)
			got, err := r.GetPayrollByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("payrollUsecase.GetPayrollByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payrollUsecase.GetPayrollByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payrollUsecase_CreatePayroll(t *testing.T) {

	now := time.Now()

	p := payslip_model.Payslip{
		ID:                    0,
		EmployeeID:            1,
		TotalAttendanceDays:   1,
		TotalOvertimeHours:    1,
		TotalAttendanceSalary: 5000,
		TotalOvertimeSalary:   1250,
		TotalReimbursement:    10000,
		TotalTakeHomePay:      16250,
		PeriodID:              1,
		CreatedAt:             now,
		UpdatedAt:             now,
		CreatedBy:             1,
		UpdatedBy:             1,
	}

	attendancePeriod := attendace_model.AttendancePeriod{
		ID:                 1,
		PeriodStart:        now.Add(-24 * time.Hour),
		PeriodEnd:          now.Add(24 * time.Hour),
		IsPayslipGenerated: false,
	}

	type args struct {
		ctx context.Context
		req dto.CreatePayrollRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(
			attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository,
			attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
			overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
			reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			employeeRepo *employee_repository_interface.MockIEmployeeRepository,
			payslipRepo *payslip_repository_interface.MockIPayslipRepository,
			transactionBundler *helper_db_test.MockITransactionBundler,
		) helper_time.TimeHelper
	}{
		{
			name: "CreatePayroll-Success",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: false,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(nil).Once()

					reimbursementRepo.On("Update", mock.Anything, &compensation_model.Reimbursement{
						UpdatedAt: now,
						UpdatedBy: 1,
						PayslipID: 1,
					}, []int{1}).Return(nil).Maybe()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "payslips" &&
							log.Action == "create" &&
							log.RecordID == 1
					})).Return(nil).Once()

					ap := &attendace_model.AttendancePeriod{
						ID:                 attendancePeriod.ID,
						IsPayslipGenerated: true,
						UpdatedAt:          now,
						UpdatedBy:          1,
					}
					attendancePeriodRepo.EXPECT().UpdatePeriod(mock.Anything, ap).Return(nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "attendances_period" &&
							log.Action == "update" &&
							log.RecordID == attendancePeriod.ID
					})).Return(nil).Once()

					_ = fn(context.Background())
				}).Return(nil).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryCreateAuditAttendancePeriod",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(nil).Once()

					reimbursementRepo.On("Update", mock.Anything, &compensation_model.Reimbursement{
						UpdatedAt: now,
						UpdatedBy: 1,
						PayslipID: 1,
					}, []int{1}).Return(nil).Maybe()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "payslips" &&
							log.Action == "create" &&
							log.RecordID == 1
					})).Return(nil).Once()

					ap := &attendace_model.AttendancePeriod{
						ID:                 attendancePeriod.ID,
						IsPayslipGenerated: true,
						UpdatedAt:          now,
						UpdatedBy:          1,
					}
					attendancePeriodRepo.EXPECT().UpdatePeriod(mock.Anything, ap).Return(nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "attendances_period" &&
							log.Action == "update" &&
							log.RecordID == attendancePeriod.ID
					})).Return(errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryUpdatePeriodAttendance",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(nil).Once()

					reimbursementRepo.On("Update", mock.Anything, &compensation_model.Reimbursement{
						UpdatedAt: now,
						UpdatedBy: 1,
						PayslipID: 1,
					}, []int{1}).Return(nil).Maybe()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "payslips" &&
							log.Action == "create" &&
							log.RecordID == 1
					})).Return(nil).Once()

					ap := &attendace_model.AttendancePeriod{
						ID:                 attendancePeriod.ID,
						IsPayslipGenerated: true,
						UpdatedAt:          now,
						UpdatedBy:          1,
					}
					attendancePeriodRepo.EXPECT().UpdatePeriod(mock.Anything, ap).Return(errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryCreateAuditPaySlip",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(nil).Once()

					reimbursementRepo.On("Update", mock.Anything, &compensation_model.Reimbursement{
						UpdatedAt: now,
						UpdatedBy: 1,
						PayslipID: 1,
					}, []int{1}).Return(nil).Maybe()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "payslips" &&
							log.Action == "create" &&
							log.RecordID == 1
					})).Return(errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryUpdateReimbursement",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(nil).Once()

					reimbursementRepo.On("Update", mock.Anything, &compensation_model.Reimbursement{
						UpdatedAt: now,
						UpdatedBy: 1,
						PayslipID: 1,
					}, []int{1}).Return(errors.New("FATAL ERROR")).Maybe()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryCreateAuditReimbursement",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, nil).Once()

					auditRepo.On("Create", mock.Anything, mock.MatchedBy(func(log audit_model.AuditLog) bool {
						return log.TableName == "reimbursements" &&
							log.Action == "update"
					})).Return(errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryCreatePayslip",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, nil).Once()

					payslipRepo.EXPECT().Create(mock.Anything, &p).Return(1, errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryEmployeeGetByID",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)

					employeeRepo.EXPECT().GetByID(mock.Anything, p.EmployeeID).Return(&employee_model.Employee{
						ID:     p.EmployeeID,
						Salary: 10000,
					}, errors.New("FATAL ERROR")).Once()

					_ = fn(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQuerySumReimbursementsByID",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().SumReimbursementsByID(mock.Anything, []int{1}).Return([]*compensation_model.EmployeeReimbursementSummary{
					{
						EmployeeID:  p.EmployeeID,
						TotalAmount: 10000,
					},
				}, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryReimbursementsGetByDateRange",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, nil).Once()

				reimbursementRepo.EXPECT().GetByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]compensation_model.Reimbursement{
					{
						ID:         1,
						Date:       now,
						Amount:     10000,
						EmployeeID: p.EmployeeID,
					},
				}, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQuerySumOvertimeByDateRange",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, nil).Once()

				overtimeRepo.EXPECT().SumOvertimeByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*compensation_model.EmployeeOvertimeSummary{
					{
						EmployeeID: p.EmployeeID,
						TotalHours: 1,
					},
				}, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryGetEmployeeCountByDateRange",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, nil).Once()
				attendanceRepo.EXPECT().GetEmployeeCountByDateRange(mock.Anything, attendancePeriod.PeriodStart, attendancePeriod.PeriodEnd).Return([]*attendace_model.EmployeeAttendanceCount{
					{
						EmployeeID: p.EmployeeID,
						Count:      1,
					},
				}, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedErrorQueryAttendancePeriodGetByID",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&attendancePeriod, errors.New("FATAL ERROR")).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedAttendancePeriodNotFound",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(nil, nil).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedAttendancePeriodAlreadyProcessed",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				a := attendancePeriod
				a.IsPayslipGenerated = true
				attendancePeriodRepo.EXPECT().GetByID(mock.Anything, p.PeriodID).Return(&a, nil).Once()

				return times
			},
		},
		{
			name: "CreatePayroll-FailedAttendancePeriodMissing",
			args: args{
				ctx: context.Background(),
				req: dto.CreatePayrollRequest{
					AttendancePeriodID: 0,
					UserID:             1,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendace_repository_interface.MockIAttendancePeriodRepository, attendanceRepo *attendace_repository_interface.MockIAttendanceRepository, overtimeRepo *compensation_repository_interface.MockIOvertimeRepository, reimbursementRepo *compensation_repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, employeeRepo *employee_repository_interface.MockIEmployeeRepository, payslipRepo *payslip_repository_interface.MockIPayslipRepository, transactionBundler *helper_db_test.MockITransactionBundler) helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				return times
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendancePeriodRepo := attendace_repository_interface.NewMockIAttendancePeriodRepository(t)
			attendanceRepo := attendace_repository_interface.NewMockIAttendanceRepository(t)
			overtimeRepo := compensation_repository_interface.NewMockIOvertimeRepository(t)
			reimbursementRepo := compensation_repository_interface.NewMockIReimbursementRepository(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			employeeRepo := employee_repository_interface.NewMockIEmployeeRepository(t)
			payslipRepo := payslip_repository_interface.NewMockIPayslipRepository(t)
			transactionBundler := helper_db_test.NewMockITransactionBundler(t)

			times := tt.beforeFunc(
				attendancePeriodRepo,
				attendanceRepo,
				overtimeRepo,
				reimbursementRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)

			r := NewPayrollUsecase(
				times,
				overtimeRepo,
				reimbursementRepo,
				attendanceRepo,
				attendancePeriodRepo,
				auditRepo,
				employeeRepo,
				payslipRepo,
				transactionBundler,
			)
			if err := r.CreatePayroll(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("payrollUsecase.CreatePayroll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
