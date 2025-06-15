package repository

import (
	"context"
	"errors"
	"time"

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

func (r *reimbursementRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]model.Reimbursement, error) {
	db := r.getDB(ctx)

	rs := []dao.ReimbursementDAO{}
	err := db.WithContext(ctx).Where("date >= ? AND date <= ?", startDate, endDate).Find(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	res := make([]model.Reimbursement, 0, len(rs))
	for _, r := range rs {
		res = append(res, model.Reimbursement{
			ID:          r.ID,
			EmployeeID:  r.EmployeeID,
			Date:        r.Date,
			Amount:      r.Amount,
			Description: r.Description,
			CreatedBy:   r.CreatedBy,
			UpdatedBy:   r.UpdatedBy,
			PayslipID:   r.PayslipID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return res, nil
}

func (r *reimbursementRepository) GetByIDs(ctx context.Context, id []int) ([]model.Reimbursement, error) {
	db := r.getDB(ctx)

	rs := []dao.ReimbursementDAO{}
	err := db.WithContext(ctx).Find(&rs, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	res := make([]model.Reimbursement, 0, len(rs))
	for _, r := range rs {
		res = append(res, model.Reimbursement{
			ID:          r.ID,
			EmployeeID:  r.EmployeeID,
			Date:        r.Date,
			Amount:      r.Amount,
			Description: r.Description,
			CreatedBy:   r.CreatedBy,
			UpdatedBy:   r.UpdatedBy,
			PayslipID:   r.PayslipID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return res, nil
}

func (r *reimbursementRepository) GetByPayslipID(ctx context.Context, payslipID int, employeeID int) ([]model.Reimbursement, error) {
	db := r.getDB(ctx)

	rs := []dao.ReimbursementDAO{}
	err := db.WithContext(ctx).Where("payslip_id = ? AND employee_id = ?", payslipID, employeeID).Find(&rs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	res := make([]model.Reimbursement, 0, len(rs))
	for _, d := range rs {
		res = append(res, model.Reimbursement{
			ID:          d.ID,
			EmployeeID:  d.EmployeeID,
			Date:        d.Date,
			Amount:      d.Amount,
			Description: d.Description,
			CreatedBy:   d.CreatedBy,
			UpdatedBy:   d.UpdatedBy,
			PayslipID:   d.PayslipID,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return res, nil
}

func (r *reimbursementRepository) SumReimbursementsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*model.EmployeeReimbursementSummary, error) {
	db := r.getDB(ctx)

	var daos []dao.ReimbursementSumDAO

	query := `
		WITH reimbursements_period AS (
			SELECT amount, employee_id
			FROM reimbursements
			WHERE date >= ? AND date <= ?
		)
		SELECT SUM(rp.amount) AS sum, rp.employee_id
		FROM reimbursements_period rp
		GROUP BY rp.employee_id;
	`

	if err := db.Raw(query, startDate, endDate).Scan(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]*model.EmployeeReimbursementSummary, 0, len(daos))
	for _, d := range daos {
		result = append(result, &model.EmployeeReimbursementSummary{
			EmployeeID:  d.EmployeeID,
			TotalAmount: d.TotalAmount,
		})
	}

	return result, nil
}

func (r *reimbursementRepository) SumReimbursementsByID(ctx context.Context, id []int) ([]*model.EmployeeReimbursementSummary, error) {
	db := r.getDB(ctx)

	var daos []dao.ReimbursementSumDAO

	query := `
		WITH reimbursements_period AS (
			SELECT amount, employee_id
			FROM reimbursements
			WHERE id IN ?
		)
		SELECT SUM(rp.amount) AS sum, rp.employee_id
		FROM reimbursements_period rp
		GROUP BY rp.employee_id;
	`

	if err := db.Raw(query, id).Scan(&daos).Error; err != nil {
		return nil, err
	}

	result := make([]*model.EmployeeReimbursementSummary, 0, len(daos))
	for _, d := range daos {
		result = append(result, &model.EmployeeReimbursementSummary{
			EmployeeID:  d.EmployeeID,
			TotalAmount: d.TotalAmount,
		})
	}

	return result, nil
}

func (r *reimbursementRepository) Update(ctx context.Context, m *model.Reimbursement, ids []int) error {
	db := r.getDB(ctx)

	updates := map[string]interface{}{
		"payslip_id": m.PayslipID,
		"updated_at": m.UpdatedAt,
		"updated_by": m.UpdatedBy,
	}

	return db.Model(&dao.ReimbursementDAO{}).
		Where("id IN ?", ids).
		Updates(updates).Error
}
