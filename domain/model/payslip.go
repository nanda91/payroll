package model

type Payslip struct {
	BaseModel
	UserID             uint    `json:"user_id"`
	PayrollPeriodID    uint    `json:"payroll_period_id"`
	BaseSalary         float64 `json:"base_salary"`
	WorkingDays        int     `json:"working_days"`
	AttendanceDays     int     `json:"attendance_days"`
	OvertimeHours      float64 `json:"overtime_hours"`
	OvertimePay        float64 `json:"overtime_pay"`
	ReimbursementTotal float64 `json:"reimbursement_total"`
	TotalPay           float64 `json:"total_pay"`

	// Relationships
	User          User          `json:"user,omitempty"`
	PayrollPeriod PayrollPeriod `json:"payroll_period,omitempty"`
}
