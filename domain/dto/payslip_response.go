package dto

import "payroll/domain/model"

type PayslipResponse struct {
	Payslip        model.Payslip         `json:"payslip"`
	Attendances    []model.Attendance    `json:"attendances"`
	Overtimes      []model.Overtime      `json:"overtimes"`
	Reimbursements []model.Reimbursement `json:"reimbursements"`
}
