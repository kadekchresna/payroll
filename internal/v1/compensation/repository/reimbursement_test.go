package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/compensation/model"
)

func Test_reimbursementRepository_Create(t *testing.T) {
	r := &model.Reimbursement{
		EmployeeID:  1,
		Date:        time.Now(),
		Amount:      100_000,
		Description: "Medical",
		PayslipID:   10,
	}
	type args struct {
		ctx context.Context
		m   *model.Reimbursement
	}
	tests := []struct {
		name       string
		args       args
		wantID     int
		wantErr    bool
		beforeFunc func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Create-Success",
			args: args{
				ctx: context.Background(),
				m:   r,
			},
			wantID:  1,
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reimbursements" ("employee_id","date","amount","description","created_at","updated_at","created_by","updated_by","payslip_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).WithArgs(r.EmployeeID, r.Date, r.Amount, r.Description, r.CreatedAt, r.UpdatedAt, r.CreatedBy, r.UpdatedBy, r.PayslipID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
		},
		{
			name: "Create-Error",
			args: args{
				ctx: context.Background(),
				m: &model.Reimbursement{
					EmployeeID:  1,
					Date:        time.Now(),
					Amount:      100_000,
					Description: "Medical",
					PayslipID:   10,
				},
			},
			wantID:  0,
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "reimbursements"`)).
					WillReturnError(errors.New("insert failed"))
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)
			r := &reimbursementRepository{db: db}

			got, err := r.Create(tt.args.ctx, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantID {
				t.Errorf("Create() got = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func Test_reimbursementRepository_GetByDateRange(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	tests := []struct {
		name       string
		want       []model.Reimbursement
		wantErr    bool
		beforeFunc func(sqlmock.Sqlmock)
	}{
		{
			name: "GetByDateRange-Success",
			want: []model.Reimbursement{
				{ID: 1, EmployeeID: 1, Amount: 10000},
			},
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE date >= $1 AND date <= $2`)).
					WithArgs(startDate, endDate).
					WillReturnRows(sqlmock.NewRows([]string{"id", "employee_id", "amount"}).AddRow(1, 1, 10000))
			},
		},
		{
			name:    "GetByDateRange-Failed",
			want:    nil,
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE date >= $1 AND date <= $2`)).
					WithArgs(startDate, endDate).
					WillReturnError(errors.New("query failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)

			r := &reimbursementRepository{db: db}
			got, err := r.GetByDateRange(context.Background(), startDate, endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetByDateRange() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_reimbursementRepository_GetByIDs(t *testing.T) {
	tests := []struct {
		name       string
		ids        []int
		want       []model.Reimbursement
		wantErr    bool
		beforeFunc func(sqlmock.Sqlmock)
	}{
		{
			name: "GetByIDs-Success",
			ids:  []int{1, 2},
			want: []model.Reimbursement{
				{ID: 1, EmployeeID: 1, Amount: 10000},
			},
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE "reimbursements"."id" IN ($1,$2)`)).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "employee_id", "amount"}).AddRow(1, 1, 10000))
			},
		},
		{
			name:    "GetByIDs-Error",
			ids:     []int{1},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE "reimbursements"."id" IN ($1)`)).
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)

			r := &reimbursementRepository{db: db}
			got, err := r.GetByIDs(context.Background(), tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetByIDs() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_reimbursementRepository_GetByPayslipID(t *testing.T) {
	tests := []struct {
		name       string
		payslipID  int
		employeeID int
		want       []model.Reimbursement
		wantErr    bool
		beforeFunc func(sqlmock.Sqlmock)
	}{
		{
			name:       "GetByPayslipID-Success",
			payslipID:  10,
			employeeID: 1,
			want: []model.Reimbursement{
				{ID: 1, EmployeeID: 1, PayslipID: 10},
			},
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE payslip_id = $1 AND employee_id = $2`)).
					WithArgs(10, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "employee_id", "payslip_id"}).AddRow(1, 1, 10))
			},
		},
		{
			name:       "GetByPayslipID-Error",
			payslipID:  99,
			employeeID: 9,
			want:       nil,
			wantErr:    true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reimbursements" WHERE payslip_id = $1 AND employee_id = $2`)).
					WithArgs(99, 9).
					WillReturnError(errors.New("query error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)

			r := &reimbursementRepository{db: db}
			got, err := r.GetByPayslipID(context.Background(), tt.payslipID, tt.employeeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByPayslipID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetByPayslipID() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_reimbursementRepository_SumReimbursementsByID(t *testing.T) {
	tests := []struct {
		name       string
		ids        []int
		want       []*model.EmployeeReimbursementSummary
		wantErr    bool
		beforeFunc func(sqlmock.Sqlmock)
	}{
		{
			name: "SumReimbursementsByID-Success",
			ids:  []int{1, 2},
			want: []*model.EmployeeReimbursementSummary{
				{EmployeeID: 1, TotalAmount: 10000},
			},
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					WITH reimbursements_period AS (
						SELECT amount, employee_id
						FROM reimbursements
						WHERE id IN ($1,$2)
					)
					SELECT SUM(rp.amount) AS sum, rp.employee_id
					FROM reimbursements_period rp
					GROUP BY rp.employee_id;
				`)).WithArgs(1, 2).WillReturnRows(
					sqlmock.NewRows([]string{"sum", "employee_id"}).AddRow(10000, 1),
				)
			},
		},
		{
			name:    "SumReimbursementsByID-Error",
			ids:     []int{1, 2},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`
					WITH reimbursements_period AS (
						SELECT amount, employee_id
						FROM reimbursements
						WHERE id IN ($1,$2)
					)
					SELECT SUM(rp.amount) AS sum, rp.employee_id
					FROM reimbursements_period rp
					GROUP BY rp.employee_id;
				`)).WithArgs(1, 2).WillReturnError(errors.New("query fail"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)

			r := &reimbursementRepository{db: db}
			got, err := r.SumReimbursementsByID(context.Background(), tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("SumReimbursementsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumReimbursementsByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reimbursementRepository_Update(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		reim       *model.Reimbursement
		ids        []int
		wantErr    bool
		beforeFunc func(sqlmock.Sqlmock)
	}{
		{
			name: "Update-Success",
			reim: &model.Reimbursement{
				PayslipID: 1,
				UpdatedAt: now,
				UpdatedBy: 1,
			},
			ids:     []int{1, 2},
			wantErr: false,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reimbursements" SET "payslip_id"=$1,"updated_at"=$2,"updated_by"=$3 WHERE id IN ($4,$5)`)).
					WithArgs(1, now, 1, 1, 2).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
		},
		{
			name: "Update-Failure",
			reim: &model.Reimbursement{
				PayslipID: 1,
				UpdatedAt: now,
				UpdatedBy: 1,
			},
			ids:     []int{1},
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reimbursements" SET`)).
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mock)

			r := &reimbursementRepository{db: db}
			err := r.Update(context.Background(), tt.reim, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
