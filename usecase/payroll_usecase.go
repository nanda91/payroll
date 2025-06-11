package usecase

import (
	"encoding/json"
	"errors"
	"payroll/domain/dto"
	"payroll/domain/model"
	"payroll/repositories"
	"time"
)

func NewPayrollUsecase(
	payrollRepo repositories.PayrollRepository,
	userRepo repositories.UserRepository,
	attendanceRepo repositories.AttendanceRepository,
	overtimeRepo repositories.OvertimeRepository,
	reimbursementRepo repositories.ReimbursementRepository,
	auditRepo repositories.AuditRepository,
) *PayrollUsecase {
	return &PayrollUsecase{
		payrollRepo:       payrollRepo,
		userRepo:          userRepo,
		attendanceRepo:    attendanceRepo,
		overtimeRepo:      overtimeRepo,
		reimbursementRepo: reimbursementRepo,
		auditRepo:         auditRepo,
	}
}

func (p *PayrollUsecase) CreatePayrollPeriod(req *dto.PayrollPeriodRequest, userID uint, ipAddress, requestID string) (*model.PayrollPeriod, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	period := &model.PayrollPeriod{
		BaseModel: model.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := p.payrollRepo.CreatePeriod(period); err != nil {
		return nil, err
	}

	// Log audit
	newData, _ := json.Marshal(period)
	p.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &userID,
		Action:    "CREATE",
		TableName: "payroll_periods",
		RecordID:  &period.ID,
		NewData:   string(newData),
	})

	return period, nil
}

func (p *PayrollUsecase) RunPayroll(req *dto.PayrollRunRequest, userID uint, ipAddress, requestID string) error {
	// Get payroll period
	period, err := p.payrollRepo.GetPeriodByID(req.PayrollPeriodID)
	if err != nil {
		return errors.New("payroll period not found")
	}

	if period.IsProcessed {
		return errors.New("payroll for this period has already been processed")
	}

	// Get all employees
	users, err := p.userRepo.GetAll()
	if err != nil {
		return err
	}

	// Process payroll for each employee
	for _, user := range users {
		if user.Role != "employee" {
			continue
		}

		payslip, err := p.calculatePayslip(&user, period)
		if err != nil {
			continue // Skip this employee if error
		}

		// Create payslip
		if err := p.payrollRepo.CreatePayslip(payslip); err != nil {
			continue // Skip if error creating payslip
		}
	}

	// Mark period as processed
	now := time.Now()
	period.IsProcessed = true
	period.ProcessedAt = &now
	period.UpdatedBy = &userID
	period.IPAddress = ipAddress
	period.RequestID = requestID

	if err := p.payrollRepo.UpdatePeriod(period); err != nil {
		return err
	}

	// Mark all records as processed
	_ = p.attendanceRepo.MarkAsProcessed(period.ID)
	_ = p.overtimeRepo.MarkAsProcessed(period.ID)
	_ = p.reimbursementRepo.MarkAsProcessed(period.ID)

	// Log audit
	newData, _ := json.Marshal(period)
	p.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &userID,
		Action:    "PAYROLL_PROCESSED",
		TableName: "payroll_periods",
		RecordID:  &period.ID,
		NewData:   string(newData),
	})

	return nil
}

func (p *PayrollUsecase) calculatePayslip(user *model.User, period *model.PayrollPeriod) (*model.Payslip, error) {
	// Calculate working days in period
	workingDays := p.calculateWorkingDays(period.StartDate, period.EndDate)

	// Get attendance records
	attendances, err := p.attendanceRepo.GetByUserAndPeriod(user.ID, period.StartDate, period.EndDate)
	if err != nil {
		return nil, err
	}

	// Get overtime records
	overtimes, err := p.overtimeRepo.GetByUserAndPeriod(user.ID, period.StartDate, period.EndDate)
	if err != nil {
		return nil, err
	}

	// Get reimbursement records
	reimbursements, err := p.reimbursementRepo.GetByUserAndPeriod(user.ID, period.StartDate, period.EndDate)
	if err != nil {
		return nil, err
	}

	// Calculate attendance days
	attendanceDays := len(attendances)

	// Calculate overtime hours
	var overtimeHours float64
	for _, overtime := range overtimes {
		overtimeHours += overtime.Hours
	}

	// Calculate reimbursement total
	var reimbursementTotal float64
	for _, reimbursement := range reimbursements {
		reimbursementTotal += reimbursement.Amount
	}

	// Calculate salary components
	dailySalary := user.Salary / 22 // Assuming 22 working days per month
	basePay := dailySalary * float64(attendanceDays)
	overtimePay := (user.Salary / 22 / 8) * 2 * overtimeHours // 2x hourly rate

	totalPay := basePay + overtimePay + reimbursementTotal

	payslip := &model.Payslip{
		BaseModel: model.BaseModel{
			CreatedBy: &user.ID,
		},
		UserID:             user.ID,
		PayrollPeriodID:    period.ID,
		BaseSalary:         user.Salary,
		WorkingDays:        workingDays,
		AttendanceDays:     attendanceDays,
		OvertimeHours:      overtimeHours,
		OvertimePay:        overtimePay,
		ReimbursementTotal: reimbursementTotal,
		TotalPay:           totalPay,
	}

	return payslip, nil
}

func (p *PayrollUsecase) calculateWorkingDays(startDate, endDate time.Time) int {
	workingDays := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workingDays++
		}
	}
	return workingDays
}

func (p *PayrollUsecase) GeneratePayslip(userID uint, periodID *uint) (*dto.PayslipResponse, error) {
	var payslip *model.Payslip
	var err error

	if periodID != nil {
		// Get specific period payslip
		payslip, err = p.payrollRepo.GetPayslipByUserAndPeriod(userID, *periodID)
		if err != nil {
			return nil, errors.New("payslip not found for specified period")
		}
	} else {
		// Get latest payslip
		payslips, err := p.payrollRepo.GetUserPayslips(userID)
		if err != nil || len(payslips) == 0 {
			return nil, errors.New("no payslip found")
		}
		payslip = &payslips[0] // Latest payslip
	}

	// Get period details
	period, err := p.payrollRepo.GetPeriodByID(payslip.PayrollPeriodID)
	if err != nil {
		return nil, err
	}

	// Get detailed records
	attendances, _ := p.attendanceRepo.GetByUserAndPeriod(userID, period.StartDate, period.EndDate)
	overtimes, _ := p.overtimeRepo.GetByUserAndPeriod(userID, period.StartDate, period.EndDate)
	reimbursements, _ := p.reimbursementRepo.GetByUserAndPeriod(userID, period.StartDate, period.EndDate)

	return &dto.PayslipResponse{
		Payslip:        *payslip,
		Attendances:    attendances,
		Overtimes:      overtimes,
		Reimbursements: reimbursements,
	}, nil
}

func (p *PayrollUsecase) GetPayrollSummary(periodID uint) (*dto.PayrollSummaryResponse, error) {
	// Get period
	period, err := p.payrollRepo.GetPeriodByID(periodID)
	if err != nil {
		return nil, errors.New("payroll period not found")
	}

	if !period.IsProcessed {
		return nil, errors.New("payroll has not been processed yet")
	}

	// Get all payslips for this period
	payslips, err := p.payrollRepo.GetPayslipsByPeriod(periodID)
	if err != nil {
		return nil, err
	}

	// Calculate total payout
	var totalPayout float64
	processedPayslips := make([]map[string]interface{}, len(payslips))
	for i, payslip := range payslips {
		userData := map[string]interface{}{
			"id":       payslip.User.ID,
			"username": payslip.User.Username,
			"role":     payslip.User.Role,
		}

		processedPayslips[i] = map[string]interface{}{
			"id":                  payslip.ID,
			"user_id":             payslip.UserID,
			"payroll_period_id":   payslip.PayrollPeriodID,
			"base_salary":         payslip.BaseSalary,
			"working_days":        payslip.WorkingDays,
			"attendance_days":     payslip.AttendanceDays,
			"overtime_hours":      payslip.OvertimeHours,
			"overtime_pay":        payslip.OvertimePay,
			"reimbursement_total": payslip.ReimbursementTotal,
			"total_pay":           payslip.TotalPay,
			"created_at":          payslip.CreatedAt,
			"updated_at":          payslip.UpdatedAt,
			"user":                userData,
		}

		totalPayout += payslip.TotalPay
	}

	return &dto.PayrollSummaryResponse{
		PayrollPeriod: *period,
		Payslips:      processedPayslips,
		TotalPayout:   totalPayout,
		EmployeeCount: len(payslips),
	}, nil
}
