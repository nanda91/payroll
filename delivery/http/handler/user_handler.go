package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payroll/domain/dto"
	"payroll/usecase"
	"payroll/utils"
)

type UserHandler struct {
	userUsecase *usecase.UserEmployeeUsecase
}

func NewUserHandler(userUsecase *usecase.UserEmployeeUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	response, err := h.userUsecase.Login(&req, ipAddress, requestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}
