package model

import "time"

type PayrollPeriod struct {
	BaseModel
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	IsProcessed bool       `gorm:"default:false" json:"is_processed"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`

	// Relationships
	Attendances    []Attendance    `json:"attendances,omitempty"`
	Overtimes      []Overtime      `json:"overtimes,omitempty"`
	Reimbursements []Reimbursement `json:"reimbursements,omitempty"`
	Payslips       []Payslip       `json:"payslips,omitempty"`
}
