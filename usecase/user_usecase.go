package usecase

import (
	"errors"
	"payroll/domain/dto"
	"payroll/domain/model"
	"payroll/repositories"
	"payroll/utils"
)

func NewUserUsecase(userRepo repositories.UserRepository, auditRepo repositories.AuditRepository) *UserEmployeeUsecase {
	return &UserEmployeeUsecase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

func (u *UserEmployeeUsecase) Login(req *dto.LoginRequest, ipAddress, requestID string) (*dto.LoginResponse, error) {
	user, err := u.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Log audit
	u.auditRepo.Create(&model.AuditLog{
		BaseModel: model.BaseModel{
			IPAddress: ipAddress,
			RequestID: requestID,
		},
		UserID:    &user.ID,
		Action:    "LOGIN",
		TableName: "users",
		RecordID:  &user.ID,
	})

	return &dto.LoginResponse{
		Token: token,
	}, nil
}

func (u *UserEmployeeUsecase) GetProfile(userID uint) (*model.User, error) {
	return u.userRepo.GetByID(userID)
}
