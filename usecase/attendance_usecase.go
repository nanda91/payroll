package usecase

import (
	"encoding/json"
	"errors"
	"payroll/domain/dto"
	"payroll/repositories"
	"time"

	"payroll/domain/model"
)

func NewAttendanceUsecase(attendanceRepo repositories.AttendanceRepository, auditRepo repositories.AuditRepository) *AttendanceUsecase {
	return &AttendanceUsecase{
		attendanceRepo: attendanceRepo,
		auditRepo:      auditRepo,
	}
}

func (a *AttendanceUsecase) SubmitAttendance(userID uint, req *dto.AttendanceRequest, ipAddress, requestID string) error {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Check if weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return errors.New("cannot submit attendance on weekends")
	}

	// Check if already submitted for this date
	existing, _ := a.attendanceRepo.GetByUserAndDate(userID, date)
	if existing != nil {
		return errors.New("attendance already submitted for this date")
	}

	// Parse check-in time
	checkIn, err := time.Parse("2006-01-02 15:04:05", req.Date+" "+req.CheckIn)
	if err != nil {
		return errors.New("invalid check-in time format")
	}

	var checkOut *time.Time
	var workingHours float64

	if req.CheckOut != "" {
		checkOutTime, err := time.Parse("2006-01-02 15:04:05", req.Date+" "+req.CheckOut)
		if err != nil {
			return errors.New("invalid check-out time format")
		}
		checkOut = &checkOutTime
		workingHours = checkOut.Sub(checkIn).Hours()
		if workingHours > 8 {
			workingHours = 8 // Cap at 8 hours for regular work
		}
	}

	attendance := &model.Attendance{
		BaseModel: model.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:       userID,
		Date:         date,
		CheckIn:      checkIn,
		CheckOut:     checkOut,
		WorkingHours: workingHours,
	}

	if err := a.attendanceRepo.Create(attendance); err != nil {
		return err
	}

	// Log audit
	newData, _ := json.Marshal(attendance)
	a.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &userID,
		Action:    "CREATE",
		TableName: "attendances",
		RecordID:  &attendance.ID,
		NewData:   string(newData),
	})

	return nil
}
