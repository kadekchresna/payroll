package dao

type ReimbursementSumDAO struct {
	EmployeeID  int     `gorm:"column:employee_id"`
	TotalAmount float64 `gorm:"column:sum"`
}
