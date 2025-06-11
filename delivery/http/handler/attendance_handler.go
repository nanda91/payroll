package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payroll/domain/dto"
	"payroll/usecase"
	"payroll/utils"
)

//type AttendanceHandler struct {
//	uc usecase.AttendanceUsecase
//}
//
//func NewAttendanceHandler(r *gin.Engine, uc usecase.AttendanceUsecase) {
//	h := &AttendanceHandler{uc: uc}
//	r.POST("/attendance", h.SubmitAttendance)
//}
//
//type submitAttendanceRequest struct {
//	EmployeeID uint      `json:"employee_id"`
//	Date       time.Time `json:"date"`
//}
//
//func (h *AttendanceHandler) SubmitAttendance(c *gin.Context) {
//	var req submitAttendanceRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := h.uc.SubmitAttendance(req.EmployeeID, req.Date)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "attendance submitted"})
//}

type AttendanceHandler struct {
	attendanceUsecase *usecase.AttendanceUsecase
}

func NewAttendanceHandler(attendanceUsecase *usecase.AttendanceUsecase) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceUsecase: attendanceUsecase,
	}
}

func (h *AttendanceHandler) SubmitAttendance(c *gin.Context) {
	var req dto.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	requestID := c.GetString("request_id")

	if err := h.attendanceUsecase.SubmitAttendance(userID, &req, ipAddress, requestID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to submit attendance", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Attendance submitted successfully", nil)
}
