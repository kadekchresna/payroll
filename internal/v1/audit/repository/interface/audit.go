package repository_interface

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/audit/model"
)

type IAuditRepository interface {
	Create(ctx context.Context, log model.AuditLog) error
}
