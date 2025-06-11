package model

type Reimbursement struct {
	BaseModel
	UserID          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	PayrollPeriodID *uint   `json:"payroll_period_id,omitempty"`
	IsProcessed     bool    `gorm:"default:false" json:"is_processed"`

	// Relationships
	User          User           `json:"user,omitempty"`
	PayrollPeriod *PayrollPeriod `json:"payroll_period,omitempty"`
}
