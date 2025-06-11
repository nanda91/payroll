package repositories

import (
	"gorm.io/gorm"
	"payroll/domain/model"
	"time"
)

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) ReimbursementRepository {
	return &reimbursementRepository{db: db}
}

func (r *reimbursementRepository) Create(reimbursement *model.Reimbursement) error {
	return r.db.Create(reimbursement).Error
}

func (r *reimbursementRepository) GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Reimbursement, error) {
	var reimbursements []model.Reimbursement
	if err := r.db.Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).Find(&reimbursements).Error; err != nil {
		return nil, err
	}
	return reimbursements, nil
}

func (r *reimbursementRepository) GetByPeriod(payrollPeriodID uint) ([]model.Reimbursement, error) {
	var reimbursements []model.Reimbursement
	if err := r.db.Where("payroll_period_id = ?", payrollPeriodID).Find(&reimbursements).Error; err != nil {
		return nil, err
	}
	return reimbursements, nil
}

func (r *reimbursementRepository) Update(reimbursement *model.Reimbursement) error {
	return r.db.Save(reimbursement).Error
}

func (r *reimbursementRepository) MarkAsProcessed(payrollPeriodID uint) error {
	return r.db.Model(&model.Reimbursement{}).
		Where("payroll_period_id = ?", payrollPeriodID).
		Update("is_processed", true).Error
}
