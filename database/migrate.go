package database

import (
	"gorm.io/gorm"
	"payroll/domain/model"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{}, &model.Attendance{}, &model.Overtime{}, &model.Reimbursement{}, &model.PayrollPeriod{}, &model.Payslip{}, &model.AuditLog{})
	if err != nil {
		return
	}
}
