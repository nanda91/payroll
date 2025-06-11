package routes

import (
	"github.com/gin-gonic/gin"
	"payroll/delivery/http/handler"
	"payroll/utils"
)

func SetupRoutes(
	userHandler *handler.UserHandler,
	attendanceHandler *handler.AttendanceHandler,
	overtimeHandler *handler.OvertimeHandler,
	reimbursementHandler *handler.ReimbursementHandler,
	payrollHandler *handler.PayrollHandler,
) *gin.Engine {
	router := gin.Default()

	//Middleware
	router.Use(utils.CORSMiddleware())
	router.Use(utils.LoggingMiddleware())
	router.Use(utils.RequestIDMiddleware())

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", userHandler.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(utils.AuthMiddleware())
	{
		// User routes
		users := api.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(utils.AdminMiddleware())
		{
			admin.POST("/payroll-periods", payrollHandler.CreatePayrollPeriod)
			admin.POST("/payroll/run", payrollHandler.RunPayroll)
			admin.GET("/payroll/summary", payrollHandler.GetPayrollSummary)
		}

		// Employee routes
		employee := api.Group("/employee")
		employee.Use(utils.EmployeeMiddleware())
		{
			employee.POST("/attendance", attendanceHandler.SubmitAttendance)
			employee.POST("/overtime", overtimeHandler.SubmitOvertime)
			employee.POST("/reimbursement", reimbursementHandler.SubmitReimbursement)
			employee.GET("/payslip", payrollHandler.GeneratePayslip)
		}
	}

	return router
}
