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
	"github.com/stretchr/testify/mock"
)

func Test_attendancePeriodUsecase_CreateAttendancePeriod(t *testing.T) {

	now := time.Now()

	req := dto.CreateAttendancePeriodRequest{
		PeriodStart: now.Add(-24 * time.Hour),
		PeriodEnd:   now.Add(24 * time.Hour),
		UserID:      1,
	}

	type args struct {
		ctx context.Context
		req *dto.CreateAttendancePeriodRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(
			attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			txBundler *helper_db_test.MockITransactionBundler,
		) *helper_time.TimeHelper
	}{
		{
			name: "CreateAttendancePeriod-Success",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			wantErr: false,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				attendancePeriodRepo.On("GetByPeriodIntersect", mock.Anything, req.PeriodStart, req.PeriodEnd).Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						ap := &model.AttendancePeriod{
							PeriodStart:        req.PeriodStart,
							PeriodEnd:          req.PeriodEnd,
							IsPayslipGenerated: false,
							CreatedAt:          now,
							UpdatedAt:          now,
							UpdatedBy:          req.UserID,
							CreatedBy:          req.UserID,
						}

						attendancePeriodRepo.On("Create", mock.Anything, ap).Return(1, nil).Once()

						oldData := map[string]interface{}{}
						newData := map[string]interface{}{
							"id":                   1,
							"period_start":         ap.PeriodStart,
							"period_end":           ap.PeriodEnd,
							"is_payslip_generated": ap.IsPayslipGenerated,
							"created_at":           ap.CreatedAt,
							"updated_at":           ap.UpdatedAt,
							"created_by":           ap.CreatedBy,
							"updated_by":           ap.UpdatedBy,
						}

						log := audit_model.AuditLog{
							TableName: "attendances_period",
							Action:    "create",
							RecordID:  1,
							OldData:   oldData,
							NewData:   newData,
							ChangedAt: now,
							ChangedBy: req.UserID,
						}

						auditRepo.On("Create", mock.Anything, log).Return(nil).Once()

						_ = fn(context.Background())

					}).Return(nil).Once()
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedErrorQueryCreateAudit",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				attendancePeriodRepo.On("GetByPeriodIntersect", mock.Anything, req.PeriodStart, req.PeriodEnd).Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						ap := &model.AttendancePeriod{
							PeriodStart:        req.PeriodStart,
							PeriodEnd:          req.PeriodEnd,
							IsPayslipGenerated: false,
							CreatedAt:          now,
							UpdatedAt:          now,
							UpdatedBy:          req.UserID,
							CreatedBy:          req.UserID,
						}

						attendancePeriodRepo.On("Create", mock.Anything, ap).Return(1, nil).Once()

						oldData := map[string]interface{}{}
						newData := map[string]interface{}{
							"id":                   1,
							"period_start":         ap.PeriodStart,
							"period_end":           ap.PeriodEnd,
							"is_payslip_generated": ap.IsPayslipGenerated,
							"created_at":           ap.CreatedAt,
							"updated_at":           ap.UpdatedAt,
							"created_by":           ap.CreatedBy,
							"updated_by":           ap.UpdatedBy,
						}

						log := audit_model.AuditLog{
							TableName: "attendances_period",
							Action:    "create",
							RecordID:  1,
							OldData:   oldData,
							NewData:   newData,
							ChangedAt: now,
							ChangedBy: req.UserID,
						}

						auditRepo.On("Create", mock.Anything, log).Return(errors.New("FATAL ERROR")).Once()

						_ = fn(context.Background())

					}).Return(errors.New("FATAL ERROR")).Once()
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedErrorQueryCreateAttendancePeriod",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				attendancePeriodRepo.On("GetByPeriodIntersect", mock.Anything, req.PeriodStart, req.PeriodEnd).Return(nil, nil).Once()

				txBundler.On("WithTransaction", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						// manually call the transaction function with a context
						fn := args.Get(1).(func(context.Context) error)

						ap := &model.AttendancePeriod{
							PeriodStart:        req.PeriodStart,
							PeriodEnd:          req.PeriodEnd,
							IsPayslipGenerated: false,
							CreatedAt:          now,
							UpdatedAt:          now,
							UpdatedBy:          req.UserID,
							CreatedBy:          req.UserID,
						}

						attendancePeriodRepo.On("Create", mock.Anything, ap).Return(1, errors.New("FATAL ERROR")).Once()

						_ = fn(context.Background())

					}).Return(errors.New("FATAL ERROR")).Once()
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedErrorQueryCreateAttendancePeriod",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				ap := &model.AttendancePeriod{
					PeriodStart:        req.PeriodStart,
					PeriodEnd:          req.PeriodEnd,
					IsPayslipGenerated: false,
					CreatedAt:          now,
					UpdatedAt:          now,
					UpdatedBy:          req.UserID,
					CreatedBy:          req.UserID,
				}
				attendancePeriodRepo.On("GetByPeriodIntersect", mock.Anything, req.PeriodStart, req.PeriodEnd).Return(ap, nil).Once()
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedErrorQueryGetByPeriodIntersect",
			args: args{
				ctx: context.Background(),
				req: &req,
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)

				attendancePeriodRepo.On("GetByPeriodIntersect", mock.Anything, req.PeriodStart, req.PeriodEnd).Return(nil, errors.New("FATAL ERROR")).Once()
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedPeriodStartEmpty",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendancePeriodRequest{
					PeriodStart: time.Time{},
					PeriodEnd:   req.PeriodEnd,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedPeriodEndEmpty",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendancePeriodRequest{
					PeriodStart: req.PeriodStart,
					PeriodEnd:   time.Time{},
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				return &times
			},
		},
		{
			name: "CreateAttendancePeriod-FailedPeriodEndBehindPeriodStart",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateAttendancePeriodRequest{
					PeriodStart: req.PeriodEnd,
					PeriodEnd:   req.PeriodStart,
				},
			},
			wantErr: true,
			beforeFunc: func(attendancePeriodRepo *attendance_repository_interface.MockIAttendancePeriodRepository, auditRepo *audit_repository_interface.MockIAuditRepository, txBundler *helper_db_test.MockITransactionBundler) *helper_time.TimeHelper {
				times := helper_time.NewTime(&now)
				return &times
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			attendancePeriodRepo := attendance_repository_interface.NewMockIAttendancePeriodRepository(t)
			txBundler := helper_db_test.NewMockITransactionBundler(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			times := tt.beforeFunc(attendancePeriodRepo, auditRepo, txBundler)
			u := NewAttendancePeriodUsecase(
				*times,
				attendancePeriodRepo,
				txBundler,
				auditRepo,
			)
			if err := u.CreateAttendancePeriod(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("attendancePeriodUsecase.CreateAttendancePeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
