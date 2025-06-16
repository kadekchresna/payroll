package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/payslip/model"
)

func Test_payslipRepository_Create(t *testing.T) {

	p := model.Payslip{
		ID:                    1,
		EmployeeID:            1,
		TotalAttendanceDays:   1,
		TotalOvertimeHours:    1,
		TotalAttendanceSalary: 1,
		TotalOvertimeSalary:   1,
		TotalReimbursement:    1,
		TotalTakeHomePay:      1,
	}

	type args struct {
		ctx context.Context
		p   *model.Payslip
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "Create-Success",
			args: args{
				ctx: context.Background(),
				p:   &p,
			},
			want:    1,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payslips" ("employee_id","total_attendance_days","total_overtime_hours","total_attendance_salary","total_overtime_salary","total_reimbursement","total_take_home_pay","created_at","updated_at","created_by","updated_by","period_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).WithArgs(p.EmployeeID, p.TotalAttendanceDays, p.TotalOvertimeHours, p.TotalAttendanceSalary, p.TotalOvertimeSalary, p.TotalReimbursement, p.TotalTakeHomePay, p.CreatedAt, p.UpdatedAt, p.CreatedBy, p.UpdatedBy, p.PeriodID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
		},
		{
			name: "Create-Failed",
			args: args{
				ctx: context.Background(),
				p:   &p,
			},
			want:    0,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "payslips" ("employee_id","total_attendance_days","total_overtime_hours","total_attendance_salary","total_overtime_salary","total_reimbursement","total_take_home_pay","created_at","updated_at","created_by","updated_by","period_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).WithArgs(p.EmployeeID, p.TotalAttendanceDays, p.TotalOvertimeHours, p.TotalAttendanceSalary, p.TotalOvertimeSalary, p.TotalReimbursement, p.TotalTakeHomePay, p.CreatedAt, p.UpdatedAt, p.CreatedBy, p.UpdatedBy, p.PeriodID).WillReturnError(errors.New("FATAL ERROR"))
				mockDB.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &payslipRepository{
				db: db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("payslipRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("payslipRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payslipRepository_GetByID(t *testing.T) {

	p := model.Payslip{
		ID:                    1,
		EmployeeID:            1,
		TotalAttendanceDays:   1,
		TotalOvertimeHours:    1,
		TotalAttendanceSalary: 1,
		TotalOvertimeSalary:   1,
		TotalReimbursement:    1,
		TotalTakeHomePay:      1,
	}

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		want       *model.Payslip
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByID-Success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &p,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."id" = $1 ORDER BY "payslips"."id" LIMIT $2`)).WithArgs(p.ID, 1).WillReturnRows(sqlmock.NewRows(
					[]string{"id", "employee_id", "total_attendance_days", "total_overtime_hours", "total_attendance_salary", "total_overtime_salary", "total_reimbursement", "total_take_home_pay", "created_at", "updated_at", "created_by", "updated_by", "period_id"},
				).AddRow(p.ID, p.EmployeeID, p.TotalAttendanceDays, p.TotalOvertimeHours, p.TotalAttendanceSalary, p.TotalOvertimeSalary, p.TotalReimbursement, p.TotalTakeHomePay, p.CreatedAt, p.UpdatedAt, p.CreatedBy, p.UpdatedBy, p.PeriodID))
			},
		},
		{
			name: "GetByID-Failed",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payslips" WHERE "payslips"."id" = $1 ORDER BY "payslips"."id" LIMIT $2`)).WithArgs(p.ID, 1).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &payslipRepository{
				db: db,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("payslipRepository.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payslipRepository.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payslipRepository_GetTotalTakeHomePayPerEmployee(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		want       []model.TotalTakeHomePayPerEmployee
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetTotalTakeHomePayPerEmployee-Success",
			args: args{
				ctx: context.Background(),
			},
			want: []model.TotalTakeHomePayPerEmployee{
				{
					EmployeeID:       1,
					TotalTakeHomePay: 10000,
				},
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT employee_id, SUM(total_take_home_pay) as total_take_home_pay FROM "payslips" GROUP BY "employee_id"`)).WillReturnRows(sqlmock.NewRows([]string{"employee_id", "total_take_home_pay"}).AddRow(1, 10000))
			},
		},
		{
			name: "GetTotalTakeHomePayPerEmployee-Failed",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT employee_id, SUM(total_take_home_pay) as total_take_home_pay FROM "payslips" GROUP BY "employee_id"`)).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &payslipRepository{
				db: db,
			}
			got, err := r.GetTotalTakeHomePayPerEmployee(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("payslipRepository.GetTotalTakeHomePayPerEmployee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payslipRepository.GetTotalTakeHomePayPerEmployee() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_payslipRepository_GetTotalTakeHomePayAllEmployees(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		want       float64
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetTotalTakeHomePayAllEmployees-Success",
			args: args{
				ctx: context.Background(),
			},
			want:    10000,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT SUM(total_take_home_pay) FROM "payslips"`)).WillReturnRows(sqlmock.NewRows([]string{"SUM(total_take_home_pay)"}).AddRow(10000))
			},
		},
		{
			name: "GetTotalTakeHomePayAllEmployees-Failed",
			args: args{
				ctx: context.Background(),
			},
			want:    0,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT SUM(total_take_home_pay) FROM "payslips"`)).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &payslipRepository{
				db: db,
			}
			got, err := r.GetTotalTakeHomePayAllEmployees(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("payslipRepository.GetTotalTakeHomePayAllEmployees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("payslipRepository.GetTotalTakeHomePayAllEmployees() = %v, want %v", got, tt.want)
			}
		})
	}
}
