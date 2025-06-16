package dao

import "time"

type UserDAO struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Username  string    `gorm:"column:username;type:varchar;not null;default:''"`
	Password  string    `gorm:"column:password;type:varchar;not null"`
	Salt      string    `gorm:"column:salt;type:varchar;not null;default:''"`
	Status    string    `gorm:"column:status;type:varchar;not null;default:'active'"`
	Role      string    `gorm:"column:role;type:varchar;not null;default:'employee'"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:false"`
	CreatedBy int       `gorm:"column:created_by;not null;default:0"`
	UpdatedBy int       `gorm:"column:updated_by;not null;default:0"`
}

func (UserDAO) TableName() string {
	return "users"
}
