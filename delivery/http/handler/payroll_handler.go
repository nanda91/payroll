package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"payroll/domain/dto"
	"payroll/usecase"
	"payroll/utils"
)

type PayrollHandler struct {
	payrollUsecase *usecase.PayrollUsecase
}

func NewPayrollHandler(payrollUsecase *usecase.PayrollUsecase) *PayrollHandler {
	return &PayrollHandler{
		payrollUsecase: payrollUsecase,
	}
}

func (h *PayrollHandler) CreatePayrollPeriod(c *gin.Context) {
	var req dto.PayrollPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	period, err := h.payrollUsecase.CreatePayrollPeriod(&req, userID, ipAddress, requestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create payroll period", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payroll period created successfully", period)
}

func (h *PayrollHandler) RunPayroll(c *gin.Context) {
	var req dto.PayrollRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	if err := h.payrollUsecase.RunPayroll(&req, userID, ipAddress, requestID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to run payroll", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payroll processed successfully", nil)
}

func (h *PayrollHandler) GeneratePayslip(c *gin.Context) {
	userID := c.GetUint("user_id")

	var periodID *uint
	if id := c.Query("period_id"); id != "" {
		var pid uint
		if n, err := fmt.Sscanf(id, "%d", &pid); err == nil && n == 1 {
			periodID = &pid
		}
	}

	payslip, err := h.payrollUsecase.GeneratePayslip(userID, periodID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Payslip not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payslip generated successfully", payslip)
}

func (h *PayrollHandler) GetPayrollSummary(c *gin.Context) {
	var periodID uint
	if n, err := fmt.Sscanf(c.Query("period_id"), "%d", &periodID); err != nil || n != 1 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid period_id", err)
		return
	}

	summary, err := h.payrollUsecase.GetPayrollSummary(periodID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to get payroll summary", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payroll summary retrieved successfully", summary)
}
