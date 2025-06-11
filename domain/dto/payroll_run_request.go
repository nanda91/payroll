package dto

type PayrollRunRequest struct {
	PayrollPeriodID uint `json:"payroll_period_id" binding:"required"`
}
