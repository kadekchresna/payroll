package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
	"github.com/kadekchresna/payroll/internal/v1/attendance/repository/dao"
	"gorm.io/gorm"
)

func Test_attendancePeriodRepository_Create(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	ap := model.AttendancePeriod{
		ID:                 1,
		PeriodStart:        now.AddDate(0, 0, -10),
		PeriodEnd:          now,
		IsPayslipGenerated: true,
	}

	type args struct {
		ctx context.Context
		ap  *model.AttendancePeriod
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
				ap:  &ap,
			},
			want:    1,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(
					`SERT INTO "attendances_period" ("is_payslip_generated","created_at","updated_at","created_by","updated_by","period_start","period_end") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "period_start","period_end","id"`)).
					WithArgs(ap.IsPayslipGenerated, ap.CreatedAt, ap.UpdatedAt, ap.CreatedBy, ap.UpdatedBy, ap.PeriodStart, ap.PeriodEnd).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
		},
		{
			name: "Create-Failed",
			args: args{
				ctx: context.Background(),
				ap:  &ap,
			},
			want:    0,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta(
					`SERT INTO "attendances_period" ("is_payslip_generated","created_at","updated_at","created_by","updated_by","period_start","period_end") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "period_start","period_end","id"`)).
					WithArgs(ap.IsPayslipGenerated, ap.CreatedAt, ap.UpdatedAt, ap.CreatedBy, ap.UpdatedBy, ap.PeriodStart, ap.PeriodEnd).
					WillReturnError(errors.New("insert failed"))
				mockDB.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &attendancePeriodRepository{
				db: db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.ap)
			if (err != nil) != tt.wantErr {
				t.Errorf("attendancePeriodRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("attendancePeriodRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attendancePeriodRepository_GetByPeriodIntersect(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	ap := dao.AttendancePeriodDAO{
		ID:                 1,
		PeriodStart:        now.AddDate(0, 0, -10),
		PeriodEnd:          now,
		IsPayslipGenerated: ptrBool(true),
	}

	type args struct {
		ctx         context.Context
		periodStart time.Time
		periodEnd   time.Time
	}
	tests := []struct {
		name       string
		args       args
		want       *model.AttendancePeriod
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByPeriodIntersect-Success",
			args: args{
				ctx:         context.Background(),
				periodStart: ap.PeriodStart,
				periodEnd:   ap.PeriodEnd,
			},
			want: &model.AttendancePeriod{
				ID:                 ap.ID,
				PeriodStart:        ap.PeriodStart,
				PeriodEnd:          ap.PeriodEnd,
				IsPayslipGenerated: true,
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendances_period" WHERE (period_start <= $1 AND period_end >= $2) OR (period_start <= $3 AND period_end >= $4) OR (period_start <= $5 AND period_end >= $6) OR (period_start >= $7 AND period_end <= $8) ORDER BY "attendances_period"."id" LIMIT $9
				`)
				mockDB.ExpectQuery(query).
					WithArgs(
						ap.PeriodEnd, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodStart,
						ap.PeriodStart, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodEnd, 1,
					).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "period_start", "period_end", "is_payslip_generated",
					}).AddRow(
						ap.ID, ap.PeriodStart, ap.PeriodEnd, *ap.IsPayslipGenerated,
					))
			},
		},
		{
			name: "GetByPeriodIntersect-NotFound",
			args: args{
				ctx:         context.Background(),
				periodStart: ap.PeriodStart,
				periodEnd:   ap.PeriodEnd,
			},
			want:    nil,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendances_period" WHERE (period_start <= $1 AND period_end >= $2) OR (period_start <= $3 AND period_end >= $4) OR (period_start <= $5 AND period_end >= $6) OR (period_start >= $7 AND period_end <= $8) ORDER BY "attendances_period"."id" LIMIT $9
				`)
				mockDB.ExpectQuery(query).
					WithArgs(
						ap.PeriodEnd, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodStart,
						ap.PeriodStart, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodEnd, 1,
					).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name: "GetByPeriodIntersect-Error",
			args: args{
				ctx:         context.Background(),
				periodStart: ap.PeriodStart,
				periodEnd:   ap.PeriodEnd,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendances_period" WHERE (period_start <= $1 AND period_end >= $2) OR (period_start <= $3 AND period_end >= $4) OR (period_start <= $5 AND period_end >= $6) OR (period_start >= $7 AND period_end <= $8) ORDER BY "attendances_period"."id" LIMIT $9
				`)
				mockDB.ExpectQuery(query).
					WithArgs(
						ap.PeriodEnd, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodStart,
						ap.PeriodStart, ap.PeriodEnd,
						ap.PeriodStart, ap.PeriodEnd, 1,
					).
					WillReturnError(errors.New("query failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &attendancePeriodRepository{
				db: db,
			}
			got, err := r.GetByPeriodIntersect(tt.args.ctx, tt.args.periodStart, tt.args.periodEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByPeriodIntersect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetByPeriodIntersect() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func ptrBool(b bool) *bool {
	return &b
}

func Test_attendancePeriodRepository_GetByID(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	ap := dao.AttendancePeriodDAO{
		ID:                 1,
		PeriodStart:        now.AddDate(0, 0, -7),
		PeriodEnd:          now,
		IsPayslipGenerated: ptrBool(true),
	}

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		want       *model.AttendancePeriod
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByID-Success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &model.AttendancePeriod{
				ID:                 ap.ID,
				PeriodStart:        ap.PeriodStart,
				PeriodEnd:          ap.PeriodEnd,
				IsPayslipGenerated: true,
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendances_period" WHERE "attendances_period"."id" = $1 ORDER BY "attendances_period"."id" LIMIT $2
				`)
				mockDB.ExpectQuery(query).
					WithArgs(ap.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "period_start", "period_end", "is_payslip_generated",
					}).AddRow(ap.ID, ap.PeriodStart, ap.PeriodEnd, *ap.IsPayslipGenerated))
			},
		},
		{
			name: "GetByID-NotFound",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want:    nil,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendances_period" WHERE "attendances_period"."id" = $1 ORDER BY "attendances_period"."id" LIMIT $2
				`)
				mockDB.ExpectQuery(query).
					WithArgs(2, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name: "GetByID-Error",
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
					SELECT * FROM "attendance_periods" WHERE "attendance_periods"."id" = $1 ORDER BY "attendance_periods"."id" LIMIT 1
				`)
				mockDB.ExpectQuery(query).
					WithArgs(3).
					WillReturnError(errors.New("query failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &attendancePeriodRepository{
				db: db,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetByID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_attendancePeriodRepository_UpdatePeriod(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	ap := model.AttendancePeriod{
		ID:                 1,
		IsPayslipGenerated: true,
		UpdatedAt:          now,
		UpdatedBy:          1,
	}

	type args struct {
		ctx context.Context
		ap  *model.AttendancePeriod
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "UpdatePeriod-Success",
			args: args{
				ctx: context.Background(),
				ap:  &ap,
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(
					`UPDATE "attendances_period" SET "is_payslip_generated"=$1,"updated_at"=$2,"updated_by"=$3 WHERE id = $4`,
				)).
					WithArgs(ap.IsPayslipGenerated, ap.UpdatedAt, ap.UpdatedBy, ap.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mockDB.ExpectCommit()
			},
		},
		{
			name: "UpdatePeriod-Failed",
			args: args{
				ctx: context.Background(),
				ap:  &ap,
			},
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()
				mockDB.ExpectExec(regexp.QuoteMeta(
					`UPDATE "attendances_period" SET "is_payslip_generated"=$1,"updated_at"=$2,"updated_by"=$3 WHERE id = $4`,
				)).
					WithArgs(ap.IsPayslipGenerated, ap.UpdatedAt, ap.UpdatedBy, ap.ID).
					WillReturnError(errors.New("update failed"))
				mockDB.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &attendancePeriodRepository{
				db: db,
			}
			err := r.UpdatePeriod(tt.args.ctx, tt.args.ap)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
