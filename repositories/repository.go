package repositories

import (
	"payroll/domain/model"
	"time"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	GetAll() ([]model.User, error)
}

type AttendanceRepository interface {
	Create(attendance *model.Attendance) error
	GetByUserAndDate(userID uint, date time.Time) (*model.Attendance, error)
	GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Attendance, error)
	GetByPeriod(payrollPeriodID uint) ([]model.Attendance, error)
	Update(attendance *model.Attendance) error
	MarkAsProcessed(payrollPeriodID uint) error
}

type OvertimeRepository interface {
	Create(overtime *model.Overtime) error
	GetByUserAndDate(userID uint, date time.Time) (*model.Overtime, error)
	GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Overtime, error)
	GetByPeriod(payrollPeriodID uint) ([]model.Overtime, error)
	Update(overtime *model.Overtime) error
	MarkAsProcessed(payrollPeriodID uint) error
}

type ReimbursementRepository interface {
	Create(reimbursement *model.Reimbursement) error
	GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Reimbursement, error)
	GetByPeriod(payrollPeriodID uint) ([]model.Reimbursement, error)
	Update(reimbursement *model.Reimbursement) error
	MarkAsProcessed(payrollPeriodID uint) error
}

type PayrollRepository interface {
	CreatePeriod(period *model.PayrollPeriod) error
	GetPeriodByID(id uint) (*model.PayrollPeriod, error)
	GetActivePeriods() ([]model.PayrollPeriod, error)
	UpdatePeriod(period *model.PayrollPeriod) error
	CreatePayslip(payslip *model.Payslip) error
	GetPayslipByUserAndPeriod(userID, periodID uint) (*model.Payslip, error)
	GetPayslipsByPeriod(periodID uint) ([]model.Payslip, error)
	GetUserPayslips(userID uint) ([]model.Payslip, error)
}

type AuditRepository interface {
	Create(log *model.AuditLog) error
	GetByUser(userID uint) ([]model.AuditLog, error)
	GetByTable(tableName string) ([]model.AuditLog, error)
}
