package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db_test "github.com/kadekchresna/payroll/infrastructure/db/helper/mocks"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendance_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface/mocks"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface/mocks"
	employee_model "github.com/kadekchresna/payroll/internal/v1/employee/model"
	employee_repository_interface "github.com/kadekchresna/payroll/internal/v1/employee/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_attendanceUsecase_CreateAttendance(t *testing.T) {
	now := time.Date(2025, 6, 16, 8, 0, 0, 0, time.UTC)
	date := now.Truncate(24 * time.Hour)

	a := model.Attendance{
		ID:          1,
		EmployeeID:  1,
		Date:        date,
		CheckedInAt: now,
	}

	type args struct {
		ctx context.Context
		req *dto.CreateAttendanceRequest
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(
			employeeRepo *employee_repository_interface.MockIEmployeeRepository,
			attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			txBundler *helper_db_test.MockITransactionBundler,
		) *attendanceUsecase
	}{
		{
			name: "Create-Success",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: false,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {

						fn := args.Get(1).(func(context.Context) error)

						attendace := a
						attendace.ID = 0
						attendace.CreatedAt = now
						attendace.UpdatedAt = now
						attendace.CreatedBy = 10
						attendace.UpdatedBy = 10

						attendanceRepo.On("Create", mock.Anything, &attendace).Return(123, nil).Once()

						oldData := map[string]interface{}{}
						newData := map[string]interface{}{
							"id":             123,
							"employee_id":    attendace.EmployeeID,
							"date":           attendace.Date,
							"created_at":     attendace.CreatedAt,
							"created_by":     attendace.CreatedBy,
							"updated_at":     attendace.UpdatedAt,
							"updated_by":     attendace.UpdatedBy,
							"checked_in_at":  attendace.CheckedInAt,
							"checked_out_at": attendace.CheckedOutAt,
						}

						log := audit_model.AuditLog{
							TableName: "attendances",
							Action:    "create",
							RecordID:  123,
							OldData:   oldData,
							NewData:   newData,
							ChangedAt: now,
							ChangedBy: 10,
						}

						auditRepo.On("Create", mock.Anything, log).Return(nil).Once()

						_ = fn(context.Background())

					}).Return(nil).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-SuccessCheckOut",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: false,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(&a, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						attendace := a
						attendace.ID = 0
						attendace.CreatedAt = now
						attendace.UpdatedAt = now
						attendace.CreatedBy = 10
						attendace.UpdatedBy = 10

						attendanceRepo.On("Create", mock.Anything, &attendace).Return(123, nil).Once()

						oldData := map[string]interface{}{
							"id":             123,
							"employee_id":    a.EmployeeID,
							"date":           a.Date,
							"created_at":     a.CreatedAt,
							"created_by":     a.CreatedBy,
							"updated_at":     a.UpdatedAt,
							"updated_by":     a.UpdatedBy,
							"checked_in_at":  a.CheckedInAt,
							"checked_out_at": a.CheckedOutAt,
						}
						newData := map[string]interface{}{
							"id":             123,
							"employee_id":    attendace.EmployeeID,
							"date":           attendace.Date,
							"created_at":     a.CreatedAt,
							"created_by":     a.CreatedBy,
							"updated_at":     attendace.UpdatedAt,
							"updated_by":     attendace.UpdatedBy,
							"checked_in_at":  attendace.CheckedInAt,
							"checked_out_at": attendace.UpdatedAt,
						}

						log := audit_model.AuditLog{
							TableName: "attendances",
							Action:    "update",
							RecordID:  123,
							OldData:   oldData,
							NewData:   newData,
							ChangedAt: now,
							ChangedBy: 10,
						}

						auditRepo.On("Create", mock.Anything, log).Return(nil).Once()

						_ = fn(context.Background())

					}).Return(nil).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedErrorQueryAudit",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						attendace := a
						attendace.ID = 0
						attendace.CreatedAt = now
						attendace.UpdatedAt = now
						attendace.CreatedBy = 10
						attendace.UpdatedBy = 10

						attendanceRepo.On("Create", mock.Anything, &attendace).Return(123, nil).Once()

						oldData := map[string]interface{}{}
						newData := map[string]interface{}{
							"id":             123,
							"employee_id":    attendace.EmployeeID,
							"date":           attendace.Date,
							"created_at":     attendace.CreatedAt,
							"created_by":     attendace.CreatedBy,
							"updated_at":     attendace.UpdatedAt,
							"updated_by":     attendace.UpdatedBy,
							"checked_in_at":  attendace.CheckedInAt,
							"checked_out_at": attendace.CheckedOutAt,
						}

						log := audit_model.AuditLog{
							TableName: "attendances",
							Action:    "create",
							RecordID:  123,
							OldData:   oldData,
							NewData:   newData,
							ChangedAt: now,
							ChangedBy: 10,
						}

						auditRepo.On("Create", mock.Anything, log).Return(errors.New("FATAL ERROR")).Once()

						_ = fn(context.Background())

					}).Return(errors.New("FATAL ERROR")).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedErrorQueryAttendance",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						attendace := a
						attendace.ID = 0
						attendace.CreatedAt = now
						attendace.UpdatedAt = now
						attendace.CreatedBy = 10
						attendace.UpdatedBy = 10

						attendanceRepo.On("Create", mock.Anything, &attendace).Return(123, errors.New("FATAL ERROR")).Once()

						_ = fn(context.Background())

					}).Return(errors.New("FATAL ERROR")).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-SuccessAlreadyCheckOut",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: false,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendace := a
				attendace.ID = 0
				attendace.CreatedAt = now
				attendace.UpdatedAt = now
				attendace.CreatedBy = 10
				attendace.UpdatedBy = 10
				attendace.CheckedOutAt = &now

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(&attendace, nil).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedErrorQueryGetByDateAndEmployeeID",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, nil).Once()

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, 1, date).
					Return(nil, errors.New("FATAL ERROR")).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedErrorQueryGetByID",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				employeeRepo.On("GetByID", mock.Anything, 1).
					Return(&employee_model.Employee{ID: 1}, errors.New("FATAL ERROR")).Once()

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedWeekend",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 1,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				now := time.Date(2025, 6, 15, 8, 0, 0, 0, time.UTC)
				times := helper_time.NewTime(&now)

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
		{
			name: "Create-FailedEmployeeIDEmpty",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendanceRequest{
					EmployeeID: 0,
					UserID:     10,
					Date:       date,
				},
			},
			wantErr: true,
			beforeFunc: func(
				employeeRepo *employee_repository_interface.MockIEmployeeRepository,
				attendanceRepo *attendance_repository_interface.MockIAttendanceRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				txBundler *helper_db_test.MockITransactionBundler,
			) *attendanceUsecase {
				times := helper_time.NewTime(&now)

				return &attendanceUsecase{
					times:              times,
					attendanceRepo:     attendanceRepo,
					employeeRepo:       employeeRepo,
					auditRepo:          auditRepo,
					transactionBundler: txBundler,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendanceRepo := attendance_repository_interface.NewMockIAttendanceRepository(t)
			employeeRepo := employee_repository_interface.NewMockIEmployeeRepository(t)
			txBundler := helper_db_test.NewMockITransactionBundler(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)

			u := tt.beforeFunc(employeeRepo, attendanceRepo, auditRepo, txBundler)

			err := u.CreateAttendance(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAttendance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
