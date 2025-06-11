package model

import "time"

type Overtime struct {
	BaseModel
	UserID          uint      `json:"user_id"`
	Date            time.Time `json:"date"`
	Hours           float64   `json:"hours"`
	Description     string    `json:"description"`
	PayrollPeriodID *uint     `json:"payroll_period_id,omitempty"`
	IsProcessed     bool      `gorm:"default:false" json:"is_processed"`

	// Relationships
	User          User           `json:"user,omitempty"`
	PayrollPeriod *PayrollPeriod `json:"payroll_period,omitempty"`
}
