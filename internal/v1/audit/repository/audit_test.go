package repository

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/audit/model"
)

func Test_auditRepo_Create(t *testing.T) {

	l := model.AuditLog{
		ID:        1,
		TableName: "table",
		Action:    "create",
		RecordID:  1,
		OldData:   map[string]interface{}{},
		NewData:   map[string]interface{}{},
		ChangedBy: 1,
	}

	type args struct {
		ctx context.Context
		log model.AuditLog
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(mockDB sqlmock.Sqlmock)
	}{
		{
			name: "Create-Success",
			args: args{
				ctx: context.Background(),
				log: l,
			},
			wantErr: false,
			beforeFunc: func(mockDB sqlmock.Sqlmock) {
				mockDB.ExpectBegin()

				oldDataJSON, _ := json.Marshal(l.OldData)
				newDataJSON, _ := json.Marshal(l.NewData)
				mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "audit_logs" ("table_name","action","record_id","old_data","new_data","changed_by","changed_at","ip_address","request_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).WithArgs(l.TableName, l.Action, l.RecordID, string(oldDataJSON), string(newDataJSON), l.ChangedBy, l.ChangedAt, l.IPAddress, l.RequestID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockDB, cleanup := helper_db.SetupMockDB(t)
			defer cleanup()

			tt.beforeFunc(mockDB)

			r := &auditRepo{
				db: db,
			}
			if err := r.Create(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("auditRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
