package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payroll/domain/dto"
	"payroll/usecase"
	"payroll/utils"
)

type OvertimeHandler struct {
	overtimeUsecase *usecase.OvertimeUsecase
}

func NewOvertimeHandler(overtimeUsecase *usecase.OvertimeUsecase) *OvertimeHandler {
	return &OvertimeHandler{
		overtimeUsecase: overtimeUsecase,
	}
}

func (h *OvertimeHandler) SubmitOvertime(c *gin.Context) {
	var req dto.OvertimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	if err := h.overtimeUsecase.SubmitOvertime(userID, &req, ipAddress, requestID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to submit overtime", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Overtime submitted successfully", nil)
}
