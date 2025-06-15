package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	audit_model "github.com/kadekchresna/payroll/internal/v1/audit/model"
	audit_repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/compensation/usecase/interface"
)

type reimbursementUsecase struct {
	times              helper_time.TimeHelper
	reimbursementRepo  repository_interface.IReimbursementRepository
	auditRepo          audit_repository_interface.IAuditRepository
	transactionBundler helper_db.ITransactionBundler
}

func NewReimbursementUsecase(
	times helper_time.TimeHelper,
	reimbursementRepo repository_interface.IReimbursementRepository,
	auditRepo audit_repository_interface.IAuditRepository,
	transactionBundler helper_db.ITransactionBundler,
) usecase_interface.IReimbursementUsecase {
	return &reimbursementUsecase{
		times:              times,
		reimbursementRepo:  reimbursementRepo,
		auditRepo:          auditRepo,
		transactionBundler: transactionBundler,
	}
}

func (u *reimbursementUsecase) CreateReimbursement(ctx context.Context, req *dto.CreateReimbursementRequest) error {

	now := u.times.Now()
	if req.Amount <= 0 {
		return errors.New("Reimbursement amount is required")
	}

	req.Description = strings.TrimSpace(req.Description)
	if len(req.Description) == 0 {
		return errors.New("Reimbursement Description is required")
	}

	err := u.transactionBundler.WithTransaction(ctx, func(ctx context.Context) error {

		r := &model.Reimbursement{
			EmployeeID:  req.EmployeeID,
			Date:        req.Date,
			Amount:      req.Amount,
			Description: req.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
			CreatedBy:   req.UserID,
			UpdatedBy:   req.UserID,
		}

		reimbursementID, err := u.reimbursementRepo.Create(ctx, r)
		if err != nil {
			e := fmt.Errorf("Failed to create reimbursement data, %s", err.Error())
			logger.LogWithContext(ctx).Error(fmt.Sprintf("error reimbursementRepo.Create. %s", e.Error()))
			return e
		}

		if err := u.auditRepo.Create(ctx, audit_model.AuditLog{
			TableName: "reimbursements",
			Action:    "create",
			RecordID:  reimbursementID,
			OldData:   map[string]interface{}{},
			NewData: map[string]interface{}{
				"id":          reimbursementID,
				"employee_id": r.EmployeeID,
				"date":        r.Date,
				"amount":      r.Amount,
				"description": r.Description,
				"created_at":  r.UpdatedAt,
				"updated_at":  r.CreatedBy,
				"created_by":  r.UpdatedBy,
				"updated_by":  r.UpdatedBy,
				"payslip_id":  r.PayslipID,
			},
			ChangedBy: r.CreatedBy,
			ChangedAt: now,
		}); err != nil {
			e := fmt.Errorf("Failed to create reimbursement audit data, %s", err.Error())
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
