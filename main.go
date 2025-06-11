package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"payroll/configs"
	"payroll/database"
	"payroll/delivery/http/handler"
	"payroll/repositories"
	"payroll/routes"
	"payroll/usecase"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}
}

func main() {
	// Initialize config
	cfg := configs.NewConfig()

	// Initialize database
	db, err := configs.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	//migrate data
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		database.Migrate(db)
		log.Println("Migration completed.")
		return
	}

	// Seed data
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		database.Seed(db)
		log.Println("Seeding completed.")
		return
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	overtimeRepo := repositories.NewOvertimeRepository(db)
	reimbursementRepo := repositories.NewReimbursementRepository(db)
	payrollRepo := repositories.NewPayrollRepository(db)
	auditRepo := repositories.NewAuditRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, auditRepo)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, auditRepo)
	overtimeUsecase := usecase.NewOvertimeUsecase(overtimeRepo, auditRepo)
	reimbursementUsecase := usecase.NewReimbursementUsecase(reimbursementRepo, auditRepo)
	payrollUsecase := usecase.NewPayrollUsecase(payrollRepo, userRepo, attendanceRepo, overtimeRepo, reimbursementRepo, auditRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUsecase)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUsecase)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementUsecase)
	payrollHandler := handler.NewPayrollHandler(payrollUsecase)

	// Setup routes
	router := routes.SetupRoutes(userHandler, attendanceHandler, overtimeHandler, reimbursementHandler, payrollHandler)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
