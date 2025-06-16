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

func Test_overtimeRepository_Create(t *testing.T) {
	now := time.Now()
	ot := model.Overtime{
		ID:         1,
		EmployeeID: 1,
		Date:       now,
		Hours:      2,
	}

	type args struct {
		ctx context.Context
		ot  *model.Overtime
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
				ot:  &ot,
			},
			want:    1,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "overtimes" ("employee_id","date","hours","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).WithArgs(ot.EmployeeID, ot.Date, ot.Hours, ot.CreatedAt, ot.UpdatedAt, ot.CreatedBy, ot.UpdatedBy).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
		},
		{
			name: "Create-Failed",
			args: args{
				ctx: context.Background(),
				ot:  &ot,
			},
			want:    0,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "overtimes" ("employee_id","date","hours","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).WithArgs(ot.EmployeeID, ot.Date, ot.Hours, ot.CreatedAt, ot.UpdatedAt, ot.CreatedBy, ot.UpdatedBy).WillReturnError(errors.New("FATAL ERROR"))
				mockDB.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &overtimeRepository{
				db: db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.ot)
			if (err != nil) != tt.wantErr {
				t.Errorf("overtimeRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("overtimeRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_overtimeRepository_GetByDateAndEmployeeID(t *testing.T) {
	date := time.Now()

	ot := model.Overtime{
		ID:         1,
		EmployeeID: 1,
		Date:       date,
		Hours:      2,
	}

	ots := []model.Overtime{ot}

	type args struct {
		ctx        context.Context
		employeeID int
		date       time.Time
	}
	tests := []struct {
		name       string
		args       args
		want       []model.Overtime
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByDateAndEmployeeID-Success",
			args: args{
				ctx:        context.Background(),
				employeeID: 1,
				date:       date,
			},
			want:    ots,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "overtimes" WHERE employee_id = $1 AND date = $2`)).WithArgs(ot.EmployeeID, ot.Date.Format(time.DateOnly)).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"employee_id",
					"date",
					"hours",
					"created_at",
					"updated_at",
					"created_by",
					"updated_by"}).AddRow(ot.ID, ot.EmployeeID, ot.Date, ot.Hours, ot.CreatedAt, ot.UpdatedAt, ot.CreatedBy, ot.UpdatedBy))
			},
		},
		{
			name: "GetByDateAndEmployeeID-Failed",
			args: args{
				ctx:        context.Background(),
				employeeID: 1,
				date:       date,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "overtimes" WHERE employee_id = $1 AND date = $2`)).WithArgs(ot.EmployeeID, ot.Date.Format(time.DateOnly)).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &overtimeRepository{
				db: db,
			}
			got, err := r.GetByDateAndEmployeeID(tt.args.ctx, tt.args.employeeID, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("overtimeRepository.GetByDateAndEmployeeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("overtimeRepository.GetByDateAndEmployeeID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_overtimeRepository_SumOvertimeByDateRange(t *testing.T) {

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)
	type args struct {
		ctx       context.Context
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name       string
		args       args
		want       []*model.EmployeeOvertimeSummary
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "SumOvertimeByDateRange-Success",
			args: args{
				ctx:       context.Background(),
				startDate: startDate,
				endDate:   endDate,
			},
			want: []*model.EmployeeOvertimeSummary{
				{
					EmployeeID: 1,
					TotalHours: 1,
				},
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
				WITH overtime_period AS (
            			SELECT hours, employee_id
            			FROM overtimes
            			WHERE date >= $1 AND date <= $2
            		)
            		SELECT SUM(op.hours) AS sum, op.employee_id
            		FROM overtime_period op
            		GROUP BY op.employee_id;
				`)).WithArgs(startDate, endDate).WillReturnRows(sqlmock.NewRows([]string{"sum", "employee_id"}).AddRow(1, 1))
			},
		},
		{
			name: "SumOvertimeByDateRange-Failed",
			args: args{
				ctx:       context.Background(),
				startDate: startDate,
				endDate:   endDate,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
				WITH overtime_period AS (
            			SELECT hours, employee_id
            			FROM overtimes
            			WHERE date >= $1 AND date <= $2
            		)
            		SELECT SUM(op.hours) AS sum, op.employee_id
            		FROM overtime_period op
            		GROUP BY op.employee_id;
				`)).WithArgs(startDate, endDate).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &overtimeRepository{
				db: db,
			}
			got, err := r.SumOvertimeByDateRange(tt.args.ctx, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("overtimeRepository.SumOvertimeByDateRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("overtimeRepository.SumOvertimeByDateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
