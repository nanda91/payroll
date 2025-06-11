package repositories

import (
	"gorm.io/gorm"
	"payroll/domain/model"
)

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (a auditRepository) Create(log *model.AuditLog) error {
	return a.db.Create(log).Error
}

func (a auditRepository) GetByUser(userID uint) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	if err := a.db.Where("user_id = ?", userID).
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (a auditRepository) GetByTable(tableName string) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	if err := a.db.Where("table_name = ?", tableName).
		Preload("User").
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
