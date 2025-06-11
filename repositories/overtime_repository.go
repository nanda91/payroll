package repositories

import (
	"gorm.io/gorm"
	"payroll/domain/model"
	"time"
)

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) OvertimeRepository {
	return &overtimeRepository{db: db}
}

func (r *overtimeRepository) Create(overtime *model.Overtime) error {
	return r.db.Create(overtime).Error
}

func (r *overtimeRepository) GetByUserAndDate(userID uint, date time.Time) (*model.Overtime, error) {
	var overtime model.Overtime
	if err := r.db.Where("user_id = ? AND DATE(date) = DATE(?)", userID, date).First(&overtime).Error; err != nil {
		return nil, err
	}
	return &overtime, nil
}

func (r *overtimeRepository) GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Overtime, error) {
	var overtimes []model.Overtime
	if err := r.db.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).Find(&overtimes).Error; err != nil {
		return nil, err
	}
	return overtimes, nil
}

func (r *overtimeRepository) GetByPeriod(payrollPeriodID uint) ([]model.Overtime, error) {
	var overtimes []model.Overtime
	if err := r.db.Where("payroll_period_id = ?", payrollPeriodID).Find(&overtimes).Error; err != nil {
		return nil, err
	}
	return overtimes, nil
}

func (r *overtimeRepository) Update(overtime *model.Overtime) error {
	return r.db.Save(overtime).Error
}

func (r *overtimeRepository) MarkAsProcessed(payrollPeriodID uint) error {
	return r.db.Model(&model.Overtime{}).
		Where("payroll_period_id = ?", payrollPeriodID).
		Update("is_processed", true).Error
}
