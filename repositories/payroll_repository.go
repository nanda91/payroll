package repositories

import (
	"gorm.io/gorm"
	"payroll/domain/model"
)

type payrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &payrollRepository{db: db}
}

func (r *payrollRepository) CreatePeriod(period *model.PayrollPeriod) error {
	return r.db.Create(period).Error
}

func (r *payrollRepository) GetPeriodByID(id uint) (*model.PayrollPeriod, error) {
	var period model.PayrollPeriod
	if err := r.db.First(&period, id).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *payrollRepository) GetActivePeriods() ([]model.PayrollPeriod, error) {
	var periods []model.PayrollPeriod
	if err := r.db.Where("is_processed = ?", false).Find(&periods).Error; err != nil {
		return nil, err
	}
	return periods, nil
}

func (r *payrollRepository) UpdatePeriod(period *model.PayrollPeriod) error {
	return r.db.Save(period).Error
}

func (r *payrollRepository) CreatePayslip(payslip *model.Payslip) error {
	return r.db.Create(payslip).Error
}

func (r *payrollRepository) GetPayslipByUserAndPeriod(userID, periodID uint) (*model.Payslip, error) {
	var payslip model.Payslip
	if err := r.db.Where("user_id = ? AND payroll_period_id = ?", userID, periodID).
		Preload("User").
		Preload("PayrollPeriod").
		First(&payslip).Error; err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (r *payrollRepository) GetPayslipsByPeriod(periodID uint) ([]model.Payslip, error) {
	var payslips []model.Payslip
	if err := r.db.Where("payroll_period_id = ?", periodID).
		Preload("User").
		Find(&payslips).Error; err != nil {
		return nil, err
	}
	return payslips, nil
}

func (r *payrollRepository) GetUserPayslips(userID uint) ([]model.Payslip, error) {
	var payslips []model.Payslip
	if err := r.db.Where("user_id = ?", userID).
		Preload("PayrollPeriod").
		Order("created_at DESC").
		Find(&payslips).Error; err != nil {
		return nil, err
	}
	return payslips, nil
}
