package model

import "time"

type Attendance struct {
	BaseModel
	UserID          uint       `json:"user_id"`
	Date            time.Time  `json:"date"`
	CheckIn         time.Time  `json:"check_in"`
	CheckOut        *time.Time `json:"check_out,omitempty"`
	WorkingHours    float64    `json:"working_hours"`
	PayrollPeriodID *uint      `json:"payroll_period_id,omitempty"`
	IsProcessed     bool       `gorm:"default:false" json:"is_processed"`

	// Relationships
	User          User           `json:"user,omitempty"`
	PayrollPeriod *PayrollPeriod `json:"payroll_period,omitempty"`
}
