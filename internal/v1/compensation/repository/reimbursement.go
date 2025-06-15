package repository

import (
	"context"

	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
	"github.com/kadekchresna/payroll/internal/v1/compensation/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/compensation/repository/interface"
	"gorm.io/gorm"
)

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) repository_interface.IReimbursementRepository {
	return &reimbursementRepository{db: db}
}

func (r *reimbursementRepository) getDB(ctx context.Context) *gorm.DB {
	if tx := helper_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *reimbursementRepository) Create(ctx context.Context, m *model.Reimbursement) (int, error) {
	db := r.getDB(ctx)
	da := dao.ReimbursementDAO{
		EmployeeID:  m.EmployeeID,
		Date:        m.Date,
		Amount:      m.Amount,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		PayslipID:   m.PayslipID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if err := db.WithContext(ctx).Create(&da).Error; err != nil {
		return 0, err
	}
	return da.ID, nil
}
