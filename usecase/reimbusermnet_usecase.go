package usecase

import (
	"encoding/json"
	"payroll/domain/dto"
	"payroll/domain/model"
	"payroll/repositories"
)

func NewReimbursementUsecase(reimbursementRepo repositories.ReimbursementRepository, auditRepo repositories.AuditRepository) *ReimbursementUsecase {
	return &ReimbursementUsecase{
		reimbursementRepo: reimbursementRepo,
		auditRepo:         auditRepo,
	}
}

func (r *ReimbursementUsecase) SubmitReimbursement(userID uint, req *dto.ReimbursementRequest, ipAddress, requestID string) error {
	reimbursement := &model.Reimbursement{
		BaseModel: model.BaseModel{
			CreatedBy: &userID,
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:      userID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	if err := r.reimbursementRepo.Create(reimbursement); err != nil {
		return err
	}

	// Log audit
	newData, _ := json.Marshal(reimbursement)
	r.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &userID,
		Action:    "CREATE",
		TableName: "reimbursements",
		RecordID:  &reimbursement.ID,
		NewData:   string(newData),
	})

	return nil
}
