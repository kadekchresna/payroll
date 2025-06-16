package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/employee/model"
)

func Test_employeeRepo_GetByID(t *testing.T) {
	e := model.Employee{
		ID:       1,
		FullName: "fullname",
		Salary:   1000000,
		Code:     "E1",
		UserID:   1,
	}

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		want       *model.Employee
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByID-Success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    &e,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."id" = $1 ORDER BY "employees"."id" LIMIT $2`)).WithArgs(e.ID, 1).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"fullname",
					"salary",
					"code",
					"user_id",
					"created_at",
					"updated_at",
					"created_by",
					"updated_by",
				}).AddRow(e.ID, e.FullName, e.Salary, e.Code, e.UserID, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy))
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
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE "employees"."id" = $1 ORDER BY "employees"."id" LIMIT $2`)).WithArgs(e.ID, 1).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &employeeRepo{
				db: db,
			}
			got, err := r.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("employeeRepo.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("employeeRepo.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_employeeRepo_GetByUserID(t *testing.T) {
	e := model.Employee{
		ID:       1,
		FullName: "fullname",
		Salary:   1000000,
		Code:     "E1",
		UserID:   1,
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name       string
		args       args
		want       *model.Employee
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "GetByUserID-Success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    &e,
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE user_id = $1 ORDER BY "employees"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"fullname",
					"salary",
					"code",
					"user_id",
					"created_at",
					"updated_at",
					"created_by",
					"updated_by",
				}).AddRow(e.ID, e.FullName, e.Salary, e.Code, e.UserID, e.CreatedAt, e.UpdatedAt, e.CreatedBy, e.UpdatedBy))
			},
		},
		{
			name: "GetByUserID-Failed",
			args: args{
				ctx:    context.Background(),
				userID: 1,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "employees" WHERE user_id = $1 ORDER BY "employees"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)
			r := &employeeRepo{
				db: db,
			}
			got, err := r.GetByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("employeeRepo.GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("employeeRepo.GetByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
