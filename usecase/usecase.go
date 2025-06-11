package usecase

import (
	"payroll/repositories"
)

type UserEmployeeUsecase struct {
	userRepo  repositories.UserRepository
	auditRepo repositories.AuditRepository
}

type AttendanceUsecase struct {
	attendanceRepo repositories.AttendanceRepository
	auditRepo      repositories.AuditRepository
}

type OvertimeUsecase struct {
	overtimeRepo repositories.OvertimeRepository
	auditRepo    repositories.AuditRepository
}

type ReimbursementUsecase struct {
	reimbursementRepo repositories.ReimbursementRepository
	auditRepo         repositories.AuditRepository
}

type PayrollUsecase struct {
	payrollRepo       repositories.PayrollRepository
	userRepo          repositories.UserRepository
	attendanceRepo    repositories.AttendanceRepository
	overtimeRepo      repositories.OvertimeRepository
	reimbursementRepo repositories.ReimbursementRepository
	auditRepo         repositories.AuditRepository
}
