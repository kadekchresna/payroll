package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/auth/model"
)

func Test_userRepo_Create(t *testing.T) {

	req := model.User{
		Username: "username",
		Password: "password",
		Salt:     "salt",
		Status:   "active",
		Role:     "admin",
	}

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(mock sqlmock.Sqlmock)
	}{
		{
			name: "CreateUser-Success",
			args: args{
				ctx:  context.Background(),
				user: &req,
			},
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password","salt","status","role","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).WithArgs(
					req.Username, req.Password, req.Salt, req.Status, req.Role, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy,
				).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "CreateUser-Failed",
			args: args{
				ctx:  context.Background(),
				user: &req,
			},
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password","salt","status","role","created_at","updated_at","created_by","updated_by") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).WithArgs(
					req.Username, req.Password, req.Salt, req.Status, req.Role, req.CreatedAt, req.UpdatedAt, req.CreatedBy, req.UpdatedBy,
				).WillReturnError(errors.New("UNEXPECTED ERROR"))

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &userRepo{
				db: db,
			}

			err := r.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.NoError(t, mockDB.ExpectationsWereMet())
			}
		})
	}
}
