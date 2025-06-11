package model

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type User struct {
	BaseModel
	Username string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Salary   float64 `gorm:"not null"`
	Role     Role    `gorm:"not null"`

	Attendances    []Attendance    `json:"attendances,omitempty"`
	Overtimes      []Overtime      `json:"overtimes,omitempty"`
	Reimbursements []Reimbursement `json:"reimbursements,omitempty"`
	Payslips       []Payslip       `json:"payslips,omitempty"`
}
