package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payroll/domain/dto"
	"payroll/usecase"
	"payroll/utils"
)

type ReimbursementHandler struct {
	reimbursementUsecase *usecase.ReimbursementUsecase
}

func NewReimbursementHandler(reimbursementUsecase *usecase.ReimbursementUsecase) *ReimbursementHandler {
	return &ReimbursementHandler{
		reimbursementUsecase: reimbursementUsecase,
	}
}

func (h *ReimbursementHandler) SubmitReimbursement(c *gin.Context) {
	var req dto.ReimbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	if err := h.reimbursementUsecase.SubmitReimbursement(userID, &req, ipAddress, requestID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to submit reimbursement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Reimbursement submitted successfully", nil)
}
