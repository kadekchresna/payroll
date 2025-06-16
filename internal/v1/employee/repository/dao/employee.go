package dao

import (
	"time"
)

type Employee struct {
	ID        int       `gorm:"primaryKey;column:id"`
	FullName  string    `gorm:"column:fullname;type:varchar;not null;default:''"`
	Salary    float64   `gorm:"column:salary;type:float8;not null;default:0"`
	Code      string    `gorm:"column:code;type:varchar;not null;default:'';index"`
	UserID    int       `gorm:"column:user_id;not null;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy int       `gorm:"column:created_by;not null;default:0"`
	UpdatedBy int       `gorm:"column:updated_by;not null;default:0"`
}

func (Employee) TableName() string {
	return "employees"
}
