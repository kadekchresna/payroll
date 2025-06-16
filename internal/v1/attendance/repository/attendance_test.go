package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	"gorm.io/gorm"
)

func Test_attendanceRepository_Create(t *testing.T) {

	a := model.Attendance{
		ID:         1,
		EmployeeID: 1,
		Date:       time.Now(),
	}

	type args struct {
		ctx context.Context
		a   *model.Attendance
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
				a:   &a,
			},
			want:    1,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendances" ("employee_id","date","checked_in_at","checked_out_at","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT ("employee_id","date") DO UPDATE SET "checked_out_at"=$9,"updated_at"=$10,"updated_by"=$11 RETURNING "id"`)).WithArgs(a.EmployeeID, a.Date, a.CheckedInAt, a.CheckedOutAt, a.CreatedAt, a.UpdatedAt, a.CreatedBy, a.UpdatedBy, a.UpdatedAt, a.UpdatedAt, a.UpdatedBy).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
		},
		{
			name: "Create-Failed",
			args: args{
				ctx: context.Background(),
				a:   &a,
			},
			want:    0,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "attendances" ("employee_id","date","checked_in_at","checked_out_at","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT ("employee_id","date") DO UPDATE SET "checked_out_at"=$9,"updated_at"=$10,"updated_by"=$11 RETURNING "id"`)).WithArgs(a.EmployeeID, a.Date, a.CheckedInAt, a.CheckedOutAt, a.CreatedAt, a.UpdatedAt, a.CreatedBy, a.UpdatedBy, a.UpdatedAt, a.UpdatedAt, a.UpdatedBy).WillReturnError(errors.New("FATAL ERROR"))
				mockDB.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &attendanceRepository{
				db: db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("attendanceRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("attendanceRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attendanceRepository_GetByDateAndEmployeeID(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	date := now.Truncate(24 * time.Hour)

	type args struct {
		ctx        context.Context
		employeeID int
		date       time.Time
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantNil    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByDateAndEmployeeID-Success",
			args: args{
				ctx:        context.Background(),
				employeeID: 1,
				date:       date,
			},
			wantErr: false,
			wantNil: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT * FROM "attendances" WHERE employee_id = $1 AND date = $2 ORDER BY "attendances"."id" LIMIT $3`)).
					WithArgs(1, date.Format("2006-01-02"), 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "employee_id", "date", "created_at", "updated_at", "created_by", "updated_by", "checked_in_at", "checked_out_at",
					}).AddRow(1, 1, date, now, now, 1, 1, now, now))
			},
		},
		{
			name: "GetByDateAndEmployeeID-FailedNotFound",
			args: args{
				ctx:        context.Background(),
				employeeID: 2,
				date:       date,
			},
			wantErr: false,
			wantNil: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT * FROM "attendances" WHERE employee_id = $1 AND date = $2 ORDER BY "attendances"."id" LIMIT $3`)).
					WithArgs(2, date.Format("2006-01-02"), 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name: "GetByDateAndEmployeeID-FailedErrorQuery",
			args: args{
				ctx:        context.Background(),
				employeeID: 3,
				date:       date,
			},
			wantErr: true,
			wantNil: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT * FROM "attendances" 
					WHERE employee_id = $1 AND date = $2 
					ORDER BY "attendances"."id" 
					LIMIT 1`)).
					WithArgs(3, date.Format("2006-01-02")).
					WillReturnError(errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &attendanceRepository{
				db: db,
			}
			got, err := r.GetByDateAndEmployeeID(tt.args.ctx, tt.args.employeeID, tt.args.date)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByDateAndEmployeeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("GetByDateAndEmployeeID() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_attendanceRepository_GetEmployeeCountByDateRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	type args struct {
		ctx         context.Context
		periodStart time.Time
		periodEnd   time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantCount  int
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetEmployeeCountByDateRange-Success",
			args: args{
				ctx:         context.Background(),
				periodStart: start,
				periodEnd:   end,
			},
			wantErr:   false,
			wantCount: 2,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					WITH attandances_period AS (
						SELECT employee_id 
						FROM attendances 
						WHERE date >= $1 AND date <= $2
					)
					SELECT employee_id, COUNT(employee_id) 
					FROM attandances_period 
					GROUP BY employee_id;
				`)).WithArgs(start, end).
					WillReturnRows(sqlmock.NewRows([]string{"employee_id", "count"}).
						AddRow(1, 10).
						AddRow(2, 5))
			},
		},
		{
			name: "GetEmployeeCountByDateRange-Failed",
			args: args{
				ctx:         context.Background(),
				periodStart: start,
				periodEnd:   end,
			},
			wantErr:   true,
			wantCount: 0,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					WITH attandances_period AS (
						SELECT employee_id 
						FROM attendances 
						WHERE date >= $1 AND date <= $2
					)
					SELECT employee_id, COUNT(employee_id) 
					FROM attandances_period 
					GROUP BY employee_id;
				`)).WithArgs(start, end).
					WillReturnError(errors.New("query error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := NewAttendanceRepository(db)
			got, err := r.GetEmployeeCountByDateRange(tt.args.ctx, tt.args.periodStart, tt.args.periodEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmployeeCountByDateRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("GetEmployeeCountByDateRange() = %v, wantCount %v", len(got), tt.wantCount)
			}
		})
	}
}
