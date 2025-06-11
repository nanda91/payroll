package usecase

import (
	"encoding/json"
	"errors"
	"payroll/domain/dto"
	"payroll/domain/model"
	"payroll/repositories"
	"time"
)

func NewOvertimeUsecase(overtimeRepo repositories.OvertimeRepository, auditRepo repositories.AuditRepository) *OvertimeUsecase {
	return &OvertimeUsecase{
		overtimeRepo: overtimeRepo,
		auditRepo:    auditRepo,
	}
}

func (o *OvertimeUsecase) SubmitOvertime(userID uint, req *dto.OvertimeRequest, ipAddress, requestID string) error {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Check if overtime already submitted for this date
	existing, _ := o.overtimeRepo.GetByUserAndDate(userID, date)
	if existing != nil {
		return errors.New("overtime already submitted for this date")
	}

	// Validate hours (max 3 hours per day)
	if req.Hours > 3 {
		return errors.New("overtime cannot exceed 3 hours per day")
	}

	overtime := &model.Overtime{
		BaseModel: model.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:      userID,
		Date:        date,
		Hours:       req.Hours,
		Description: req.Description,
	}

	if err := o.overtimeRepo.Create(overtime); err != nil {
		return err
	}

	// Log audit
	newData, _ := json.Marshal(overtime)
	err = o.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &userID,
		Action:    "CREATE",
		TableName: "overtimes",
		RecordID:  &overtime.ID,
		NewData:   string(newData),
	})
	if err != nil {
		return err
	}

	return nil
}
