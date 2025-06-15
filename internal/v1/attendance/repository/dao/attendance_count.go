package dao

type AttendanceCountDAO struct {
	EmployeeID int `gorm:"column:employee_id"`
	Count      int `gorm:"column:count"`
}
