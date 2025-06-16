package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper/mocks"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface/mocks"
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_reimbursementUsecase_CreateReimbursement(t *testing.T) {

	now := time.Now()

	r := model.Reimbursement{
		EmployeeID:  1,
		Date:        now,
		Amount:      100_000,
		Description: "Medical",
		PayslipID:   0,
	}

	type args struct {
		ctx context.Context
		req *dto.CreateReimbursementRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(
			reimbursementRepo *repository_interface.MockIReimbursementRepository,
			auditRepo *audit_repository_interface.MockIAuditRepository,
			transactionBundler *helper_db.MockITransactionBundler,
		) *reimbursementUsecase
	}{
		{
			name: "CreateReimbursement-Success",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateReimbursementRequest{
					EmployeeID:  r.EmployeeID,
					UserID:      r.CreatedBy,
					Amount:      r.Amount,
					Description: r.Description,
					Date:        r.Date,
				},
			},
			wantErr: false,
			beforeFunc: func(reimbursementRepo *repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, transactionBundler *helper_db.MockITransactionBundler) *reimbursementUsecase {

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

					txFunc := args.Get(1).(func(context.Context) error)
					d := r
					d.CreatedAt = now
					d.UpdatedAt = now
					reimbursementRepo.On("Create", mock.Anything, &d).
						Return(1, nil).Once()

					log := audit_model.AuditLog{
						TableName: "reimbursements",
						Action:    "create",
						RecordID:  1,
						OldData:   map[string]interface{}{},
						NewData: map[string]interface{}{
							"id":          1,
							"employee_id": d.EmployeeID,
							"date":        d.Date,
							"amount":      d.Amount,
							"description": d.Description,
							"created_at":  d.UpdatedAt,
							"updated_at":  d.CreatedBy,
							"created_by":  d.UpdatedBy,
							"updated_by":  d.UpdatedBy,
							"payslip_id":  d.PayslipID,
						},
						ChangedBy: r.CreatedBy,
						ChangedAt: now,
					}
					auditRepo.On("Create", mock.Anything, log).
						Return(nil).Once()
					_ = txFunc(context.Background())
				}).Return(nil).Once()

				return &reimbursementUsecase{
					times:              helper_time.NewTime(&now),
					reimbursementRepo:  reimbursementRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateReimbursement-FailedErrorQueryCreateAudit",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateReimbursementRequest{
					EmployeeID:  r.EmployeeID,
					UserID:      r.CreatedBy,
					Amount:      r.Amount,
					Description: r.Description,
					Date:        r.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(reimbursementRepo *repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, transactionBundler *helper_db.MockITransactionBundler) *reimbursementUsecase {

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

					txFunc := args.Get(1).(func(context.Context) error)
					d := r
					d.CreatedAt = now
					d.UpdatedAt = now
					reimbursementRepo.On("Create", mock.Anything, &d).
						Return(1, nil).Once()

					log := audit_model.AuditLog{
						TableName: "reimbursements",
						Action:    "create",
						RecordID:  1,
						OldData:   map[string]interface{}{},
						NewData: map[string]interface{}{
							"id":          1,
							"employee_id": d.EmployeeID,
							"date":        d.Date,
							"amount":      d.Amount,
							"description": d.Description,
							"created_at":  d.UpdatedAt,
							"updated_at":  d.CreatedBy,
							"created_by":  d.UpdatedBy,
							"updated_by":  d.UpdatedBy,
							"payslip_id":  d.PayslipID,
						},
						ChangedBy: r.CreatedBy,
						ChangedAt: now,
					}
					auditRepo.On("Create", mock.Anything, log).
						Return(errors.New("FATAL ERROR")).Once()
					_ = txFunc(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return &reimbursementUsecase{
					times:              helper_time.NewTime(&now),
					reimbursementRepo:  reimbursementRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateReimbursement-FailedErrorQueryCreateReimbursement",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateReimbursementRequest{
					EmployeeID:  r.EmployeeID,
					UserID:      r.CreatedBy,
					Amount:      r.Amount,
					Description: r.Description,
					Date:        r.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(reimbursementRepo *repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, transactionBundler *helper_db.MockITransactionBundler) *reimbursementUsecase {

				transactionBundler.On("WithTransaction", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {

					txFunc := args.Get(1).(func(context.Context) error)
					d := r
					d.CreatedAt = now
					d.UpdatedAt = now
					reimbursementRepo.On("Create", mock.Anything, &d).
						Return(1, errors.New("FATAL ERROR")).Once()

					_ = txFunc(context.Background())
				}).Return(errors.New("FATAL ERROR")).Once()

				return &reimbursementUsecase{
					times:              helper_time.NewTime(&now),
					reimbursementRepo:  reimbursementRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateReimbursement-FailedAmountEmpty",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateReimbursementRequest{
					EmployeeID:  r.EmployeeID,
					UserID:      r.CreatedBy,
					Amount:      0,
					Description: r.Description,
					Date:        r.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(reimbursementRepo *repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, transactionBundler *helper_db.MockITransactionBundler) *reimbursementUsecase {

				return &reimbursementUsecase{
					times:              helper_time.NewTime(&now),
					reimbursementRepo:  reimbursementRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
		{
			name: "CreateReimbursement-FailedDescEmpty",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateReimbursementRequest{
					EmployeeID:  r.EmployeeID,
					UserID:      r.CreatedBy,
					Amount:      r.Amount,
					Description: "",
					Date:        r.Date,
				},
			},
			wantErr: true,
			beforeFunc: func(reimbursementRepo *repository_interface.MockIReimbursementRepository, auditRepo *audit_repository_interface.MockIAuditRepository, transactionBundler *helper_db.MockITransactionBundler) *reimbursementUsecase {

				return &reimbursementUsecase{
					times:              helper_time.NewTime(&now),
					reimbursementRepo:  reimbursementRepo,
					auditRepo:          auditRepo,
					transactionBundler: transactionBundler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reimbursementRepo := repository_interface.NewMockIReimbursementRepository(t)
			auditRepo := audit_repository_interface.NewMockIAuditRepository(t)
			transactionBundler := helper_db.NewMockITransactionBundler(t)

			u := tt.beforeFunc(reimbursementRepo, auditRepo, transactionBundler)
			if err := u.CreateReimbursement(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("reimbursementUsecase.CreateReimbursement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
