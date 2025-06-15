package model

import "time"

type AuditLog struct {
	ID        int                    `json:"id"`
	TableName string                 `json:"table_name"`
	Action    string                 `json:"action"`
	RecordID  int                    `json:"record_id"`
	OldData   map[string]interface{} `json:"old_data,omitempty"`
	NewData   map[string]interface{} `json:"new_data,omitempty"`
	ChangedBy int                    `json:"changed_by"`
	ChangedAt time.Time              `json:"changed_at"`
	IPAddress string                 `json:"ip_address"`
	RequestID string                 `json:"request_id"`
}
