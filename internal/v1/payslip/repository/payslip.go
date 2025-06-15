package repository

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/payslip/model"
	"github.com/kadekchresna/payroll/internal/v1/payslip/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/payslip/repository/interface"
	"gorm.io/gorm"
)

type payslipRepository struct {
	db *gorm.DB
}

func NewPayslipRepository(db *gorm.DB) repository_interface.IPayslipRepository {
	return &payslipRepository{db: db}
}

func (r *payslipRepository) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *payslipRepository) Create(ctx context.Context, p *model.Payslip) (int, error) {
	db := r.getDB(ctx)

	daoModel := dao.PayslipDAO{
		EmployeeID:            p.EmployeeID,
		TotalAttendanceDays:   p.TotalAttendanceDays,
		TotalOvertimeHours:    p.TotalOvertimeHours,
		TotalAttendanceSalary: p.TotalAttendanceSalary,
		TotalOvertimeSalary:   p.TotalOvertimeSalary,
		TotalReimbursement:    p.TotalReimbursement,
		TotalTakeHomePay:      p.TotalTakeHomePay,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
		CreatedBy:             p.CreatedBy,
		UpdatedBy:             p.UpdatedBy,
		PeriodID:              p.PeriodID,
	}

	if err := db.WithContext(ctx).Create(&daoModel).Error; err != nil {
		return 0, err
	}

	return daoModel.ID, nil
}

func (r *payslipRepository) GetByID(ctx context.Context, id int) (*model.Payslip, error) {
	db := r.getDB(ctx)

	var daoModel dao.PayslipDAO
	if err := db.WithContext(ctx).First(&daoModel, id).Error; err != nil {
		return nil, err
	}

	p := &model.Payslip{
		ID:                    daoModel.ID,
		EmployeeID:            daoModel.EmployeeID,
		TotalAttendanceDays:   daoModel.TotalAttendanceDays,
		TotalOvertimeHours:    daoModel.TotalOvertimeHours,
		TotalAttendanceSalary: daoModel.TotalAttendanceSalary,
		TotalOvertimeSalary:   daoModel.TotalOvertimeSalary,
		TotalReimbursement:    daoModel.TotalReimbursement,
		TotalTakeHomePay:      daoModel.TotalTakeHomePay,
		CreatedAt:             daoModel.CreatedAt,
		UpdatedAt:             daoModel.UpdatedAt,
		CreatedBy:             daoModel.CreatedBy,
		UpdatedBy:             daoModel.UpdatedBy,
		PeriodID:              daoModel.PeriodID,
	}

	return p, nil
}

func (r *payslipRepository) GetTotalTakeHomePayPerEmployee(ctx context.Context) ([]model.TotalTakeHomePayPerEmployee, error) {
	var data []dao.TotalTakeHomePayPerEmployee

	err := r.db.WithContext(ctx).
		Table("payslips").
		Select("employee_id, SUM(total_take_home_pay) as total_take_home_pay").
		Group("employee_id").
		Scan(&data).Error

	if err != nil {
		return nil, err
	}

	res := make([]model.TotalTakeHomePayPerEmployee, 0, len(data))
	for _, d := range data {
		res = append(res, model.TotalTakeHomePayPerEmployee{
			EmployeeID:       d.EmployeeID,
			TotalTakeHomePay: d.TotalTakeHomePay,
		})
	}

	return res, nil
}

func (r *payslipRepository) GetTotalTakeHomePayAllEmployees(ctx context.Context) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Table("payslips").
		Select("SUM(total_take_home_pay)").
		Scan(&total).Error

	if err != nil {
		return 0, err
	}

	return total, nil
}
