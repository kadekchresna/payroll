package dao

type OvertimeSumDAO struct {
	EmployeeID int `gorm:"column:employee_id"`
	TotalHours int `gorm:"column:sum"`
}
