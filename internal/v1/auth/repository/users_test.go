package repository

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/auth/model"
	"github.com/stretchr/testify/assert"
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

func Test_userRepo_GetByUsername(t *testing.T) {

	res := model.User{
		ID:       1,
		Username: "username",
	}

	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name       string
		args       args
		want       *model.User
		wantErr    bool
		beforeFunc func(mock sqlmock.Sqlmock)
	}{
		{
			name: "GetByUsername-Success",
			args: args{
				ctx:      context.Background(),
				username: res.Username,
			},
			want: &res,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(res.Username, 1).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"username",
					"password",
					"salt",
					"status",
					"role",
					"created_at",
					"updated_at",
					"created_by",
					"updated_by",
				}).AddRow(res.ID, res.Username, res.Password, res.Salt, res.Status, res.Role, res.CreatedAt, res.UpdatedAt, res.CreatedBy, res.UpdatedBy))
			},
		},
		{
			name: "GetByUsername-Failed",
			args: args{
				ctx:      context.Background(),
				username: res.Username,
			},
			want:    nil,
			wantErr: true,
			beforeFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(res.Username, 1).WillReturnError(errors.New("FATAL ERROR"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := NewUserRepo(db)
			got, err := r.GetByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepo.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.GetByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
