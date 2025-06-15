package dto

type GetPayrollSummary struct {
	EmployeePayrollSummary []EmployeePayroll `json:"employee_payroll_summary"`
	TotalTakeHomePayAll    float64           `json:"total_take_home_pay_all"`
}

type EmployeePayroll struct {
	EmployeeID       int
	TotalTakeHomePay float64
}
