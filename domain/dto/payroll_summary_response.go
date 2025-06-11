package dto

import "payroll/domain/model"

type PayrollSummaryResponse struct {
	PayrollPeriod model.PayrollPeriod      `json:"payroll_period"`
	Payslips      []map[string]interface{} `json:"payslips"`
	TotalPayout   float64                  `json:"total_payout"`
	EmployeeCount int                      `json:"employee_count"`
}
