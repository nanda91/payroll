package repositories

import (
	"time"

	"gorm.io/gorm"
	"payroll/domain/model"
)

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(attendance *model.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) GetByUserAndDate(userID uint, date time.Time) (*model.Attendance, error) {
	var attendance model.Attendance
	if err := r.db.Where("user_id = ? AND DATE(date) = DATE(?)", userID, date).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) GetByUserAndPeriod(userID uint, startDate, endDate time.Time) ([]model.Attendance, error) {
	var attendances []model.Attendance
	if err := r.db.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *attendanceRepository) GetByPeriod(payrollPeriodID uint) ([]model.Attendance, error) {
	var attendances []model.Attendance
	if err := r.db.Where("payroll_period_id = ?", payrollPeriodID).Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *attendanceRepository) Update(attendance *model.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *attendanceRepository) MarkAsProcessed(payrollPeriodID uint) error {
	return r.db.Model(&model.Attendance{}).
		Where("payroll_period_id = ?", payrollPeriodID).
		Update("is_processed", true).Error
}
