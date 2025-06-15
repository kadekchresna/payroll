package repository

import (
	"context"
	"encoding/json"

	"github.com/kadekchresna/payroll/helper/logger"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	"github.com/kadekchresna/payroll/internal/v1/audit/model"
	"github.com/kadekchresna/payroll/internal/v1/audit/repository/dao"
	repository_interface "github.com/kadekchresna/payroll/internal/v1/audit/repository/interface"
	"gorm.io/gorm"
)

type auditRepo struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) repository_interface.IAuditRepository {
	return &auditRepo{db: db}
}

func (r *auditRepo) getDB(ctx context.Context) *gorm.DB {
	if tx := helper_db.GetTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *auditRepo) Create(ctx context.Context, log model.AuditLog) error {

	db := r.getDB(ctx)
	oldDataJSON, _ := json.Marshal(log.OldData)
	newDataJSON, _ := json.Marshal(log.NewData)

	requestID, _ := ctx.Value(logger.RequestIDKey).(string)
	clientIP, _ := ctx.Value(logger.ClientIPKey).(string)

	return db.WithContext(ctx).Create(&dao.AuditLog{
		TableNames: log.TableName,
		Action:     log.Action,
		RecordID:   log.RecordID,
		OldData:    oldDataJSON,
		NewData:    newDataJSON,
		ChangedBy:  log.ChangedBy,
		IPAddress:  clientIP,
		RequestID:  requestID,
	}).Error
}
