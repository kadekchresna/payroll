package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper/mocks"
	attendance_model "github.com/kadekchresna/payroll/internal/v1/attendance/model"
	attendace_repository_interface "github.com/kadekchresna/payroll/internal/v1/attendance/repository/interface/mocks"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface/mocks"
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
	compensation_model "github.com/kadekchresna/payroll/internal/v1/compensation/model"
	compensation_repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_overtimeUsecase_CreateOvertime(t *testing.T) {

	now := time.Now()
	checkedOutAt := now

	overtimeModel := compensation_model.Overtime{
		EmployeeID: 1,
		Date:       now,
		Hours:      2,
		CreatedBy:  100,
		UpdatedBy:  100,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	type args struct {
		ctx context.Context
		req *dto.CreateOvertimeRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(
			attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
			overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			transactionBundler *helper_db.MockITransactionBundler,
		) *overtimeUsecase
	}{
		{
			name: "CreateOvertime-Success",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: false,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: &checkedOutAt,
					}, nil).Once()

				overtimeRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return([]compensation_model.Overtime{}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						txFunc := args.Get(1).(func(context.Context) error)

						overtimeRepo.On("Create", mock.Anything, &overtimeModel).
							Return(1, nil).Once()

						a := audit_model.AuditLog{
							TableName: "overtimes",
							Action:    "create",
							RecordID:  1,
							OldData:   map[string]interface{}{},
							NewData: map[string]interface{}{
								"id":          1,
								"employee_id": overtimeModel.EmployeeID,
								"date":        overtimeModel.Date,
								"hours":       overtimeModel.Hours,
								"created_at":  overtimeModel.CreatedAt,
								"updated_at":  overtimeModel.UpdatedAt,
								"created_by":  overtimeModel.CreatedBy,
								"updated_by":  overtimeModel.UpdatedBy,
							},
							ChangedBy: overtimeModel.CreatedBy,
							ChangedAt: now,
						}
						auditRepo.On("Create", mock.Anything, a).Return(nil).Once()

						_ = txFunc(context.Background())
					}).Return(nil).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedErrorQueryCreateAudit",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: &checkedOutAt,
					}, nil).Once()

				overtimeRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return([]compensation_model.Overtime{}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						txFunc := args.Get(1).(func(context.Context) error)

						overtimeRepo.On("Create", mock.Anything, &overtimeModel).
							Return(1, nil).Once()

						a := audit_model.AuditLog{
							TableName: "overtimes",
							Action:    "create",
							RecordID:  1,
							OldData:   map[string]interface{}{},
							NewData: map[string]interface{}{
								"id":          1,
								"employee_id": overtimeModel.EmployeeID,
								"date":        overtimeModel.Date,
								"hours":       overtimeModel.Hours,
								"created_at":  overtimeModel.CreatedAt,
								"updated_at":  overtimeModel.UpdatedAt,
								"created_by":  overtimeModel.CreatedBy,
								"updated_by":  overtimeModel.UpdatedBy,
							},
							ChangedBy: overtimeModel.CreatedBy,
							ChangedAt: now,
						}
						auditRepo.On("Create", mock.Anything, a).Return(errors.New("FATAL ERROR")).Once()

						_ = txFunc(context.Background())
					}).Return(errors.New("FATAL ERROR")).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedErrorQueryGetByDateAndEmployeeID",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: &checkedOutAt,
					}, nil).Once()

				overtimeRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return([]compensation_model.Overtime{}, errors.New("FATAL ERROR")).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedErrorQueryCreateOvertime",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: &checkedOutAt,
					}, nil).Once()

				overtimeRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return([]compensation_model.Overtime{}, nil).Once()

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						txFunc := args.Get(1).(func(context.Context) error)

						overtimeRepo.On("Create", mock.Anything, &overtimeModel).
							Return(1, errors.New("FATAL ERROR")).Once()

						_ = txFunc(context.Background())
					}).Return(errors.New("FATAL ERROR")).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-Failed-HoursAboveLimit",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      4,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: &checkedOutAt,
					}, nil).Once()

				overtimeRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return([]compensation_model.Overtime{
						{Hours: 2},
					}, nil).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedEmployeeIDMissing",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: 0,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedHoursIDMissing",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      0,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedNotCheckedOutYet",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: nil,
					}, nil).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateOvertime-FailedErrorQueryGetByDateAndEmployeeID",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateOvertimeRequest{
					EmployeeID: overtimeModel.EmployeeID,
					UserID:     overtimeModel.CreatedBy,
					Hours:      overtimeModel.Hours,
					Date:       overtimeModel.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(attendanceRepo *attendace_repository_interface.MockIAttendanceRepository,
				overtimeRepo *compensation_repository_interface.MockIOvertimeRepository,
				auditRepo *audit_repository_interface.MockIAuditRepository,
				transactionBundler *helper_db.MockITransactionBundler) *overtimeUsecase {

				attendanceRepo.On("GetByDateAndEmployeeID", mock.Anything, overtimeModel.EmployeeID, overtimeModel.Date).
					Return(&attendance_model.Attendance{
						CheckedOutAt: nil,
					}, errors.New("FATAL ERROR")).Once()

				return &overtimeUsecase{
					times:              helper_time.NewTime(&now),
					attendanceRepo:     attendanceRepo,
					overtimeRepo:       overtimeRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendanceRepo := attendace_repository_interface.NewMockIAttendanceRepository(t)
			overtimeRepo := compensation_repository_interface.NewMockIOvertimeRepository(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			transactionBundler := helper_db.NewMockITransactionBundler(t)

			u := tt.beforeFunc(attendanceRepo, overtimeRepo, auditRepo, transactionBundler)
			if err := u.CreateOvertime(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("overtimeUsecase.CreateOvertime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
