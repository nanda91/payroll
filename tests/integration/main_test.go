// tests/integration/main_test.go
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"payroll/routes"
	"payroll/utils"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"payroll/delivery/http/handler"
	"payroll/domain/model"
	"payroll/repositories"
	"payroll/usecase"
)

type TestSuite struct {
	suite.Suite
	db            *gorm.DB
	router        *gin.Engine
	container     *pg.PostgresContainer
	adminToken    string
	employeeToken string
	adminUser     *model.User
	employeeUser  *model.User
}

func (s *TestSuite) SetupSuite() {
	if runtime.GOOS == "darwin" {
		os.Setenv("TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE", "/var/run/docker.sock")
		os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
	}
	ctx := context.Background()

	// 1. Start container dengan health check
	pgContainer, err := pg.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		pg.WithDatabase("payroll_test"),
		pg.WithUsername("tests"),
		pg.WithPassword("tests"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(s.T(), err)
	s.container = pgContainer

	// 2. Get mapped port
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(s.T(), err, "Failed to get mapped port")

	// 3. Build connection string
	connStr := fmt.Sprintf("host=localhost port=%s user=tests dbname=payroll_test password=tests sslmode=disable connect_timeout=5", port.Port())

	// 4. Connect with retry
	var db *gorm.DB
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
			}
			if err == nil {
				s.db = db
				break
			}
		}

		if i == maxAttempts-1 {
			require.NoError(s.T(), err, "Failed to connect after %d attempts", maxAttempts)
		}
		time.Sleep(2 * time.Second)
	}

	// 5. Run migrations
	s.runMigrations()
	s.setupTestData()
	s.setupRouter()
}

func (s *TestSuite) TearDownSuite() {
	if s.container != nil {
		ctx := context.Background()
		if err := s.container.Terminate(ctx); err != nil {
			s.T().Logf("Warning: failed to terminate container: %v", err)
		}
	}
}

func (s *TestSuite) SetupTest() {
	// Start a transaction for each tests
	tx := s.db.Begin()
	s.db = tx
}

func (s *TestSuite) runMigrations() {
	models := []interface{}{
		&model.User{},
		&model.Attendance{},
		&model.Overtime{},
		&model.Reimbursement{},
		&model.PayrollPeriod{},
		&model.Payslip{},
		&model.AuditLog{},
	}

	for _, model := range models {
		err := s.db.AutoMigrate(model)
		require.NoError(s.T(), err, "Failed to migrate model %T", model)
	}
}

func (s *TestSuite) setupTestData() {
	// Create admin user
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	adminUser := &model.User{
		Username: "admin",
		Role:     "admin",
		Password: string(password),
	}
	result := s.db.Create(adminUser)
	require.NoError(s.T(), result.Error, "Failed to create admin user")
	s.adminUser = adminUser

	// Create employee user
	employeeUser := &model.User{
		Username: "employee",
		Role:     "employee",
		Password: string(password),
	}
	result = s.db.Create(employeeUser)
	require.NoError(s.T(), result.Error, "Failed to create employee user")
	s.employeeUser = employeeUser

	// Generate tests tokens
	adminToken, _ := utils.GenerateToken(adminUser)
	employeeToken, _ := utils.GenerateToken(employeeUser)
	s.adminToken = fmt.Sprintf("Bearer %s", adminToken)
	s.employeeToken = fmt.Sprintf("Bearer %s", employeeToken)
}

func (s *TestSuite) setupRouter() {
	gin.SetMode(gin.TestMode)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(s.db)
	attendanceRepo := repositories.NewAttendanceRepository(s.db)
	overtimeRepo := repositories.NewOvertimeRepository(s.db)
	reimbursementRepo := repositories.NewReimbursementRepository(s.db)
	payrollRepo := repositories.NewPayrollRepository(s.db)
	auditRepo := repositories.NewAuditRepository(s.db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, auditRepo)
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepo, auditRepo)
	overtimeUsecase := usecase.NewOvertimeUsecase(overtimeRepo, auditRepo)
	reimbursementUsecase := usecase.NewReimbursementUsecase(reimbursementRepo, auditRepo)
	payrollUsecase := usecase.NewPayrollUsecase(
		payrollRepo, userRepo, attendanceRepo,
		overtimeRepo, reimbursementRepo, auditRepo,
	)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	attendanceHandler := handler.NewAttendanceHandler(attendanceUsecase)
	overtimeHandler := handler.NewOvertimeHandler(overtimeUsecase)
	reimbursementHandler := handler.NewReimbursementHandler(reimbursementUsecase)
	payrollHandler := handler.NewPayrollHandler(payrollUsecase)

	// Setup routes
	s.router = routes.SetupRoutes(
		userHandler,
		attendanceHandler,
		overtimeHandler,
		reimbursementHandler,
		payrollHandler,
	)
}

func (s *TestSuite) makeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(s.T(), err, "Failed to marshal request body")
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

// Authentication Tests
func (s *TestSuite) TestAuthentication() {
	s.Run("Successful Admin Login", func() {
		loginData := map[string]string{
			"username": "admin",
			"password": "password",
		}

		w := s.makeRequest("POST", "/api/auth/login", loginData, "")
		assert.Equal(s.T(), http.StatusOK, w.Code, "Expected status 200 for successful login")
	})

	s.Run("Invalid Credentials", func() {
		loginData := map[string]string{
			"username": "admin",
			"password": "wrongpassword",
		}

		w := s.makeRequest("POST", "/api/auth/login", loginData, "")
		assert.Equal(s.T(), http.StatusUnauthorized, w.Code, "Expected status 401 for invalid credentials")
	})
}

func (s *TestSuite) createTestPayrollPeriod() *model.PayrollPeriod {
	now := time.Now()
	period := &model.PayrollPeriod{
		StartDate:   now.AddDate(0, -1, 0),
		EndDate:     now.AddDate(0, 0, -1),
		IsProcessed: false,
	}
	result := s.db.Create(period)
	require.NoError(s.T(), result.Error, "Failed to create tests payroll period")
	return period
}

func (s *TestSuite) createTestAttendance() *model.Attendance {
	now := time.Now()
	checkout := now.Add(8 * time.Hour)
	attendance := &model.Attendance{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedBy: &s.employeeUser.ID,
			IPAddress: "127.0.0.1",
			RequestID: "tests-request-id",
		},
		UserID:   s.employeeUser.ID,
		Date:     now,
		CheckIn:  now,
		CheckOut: &checkout,
	}
	result := s.db.Create(attendance)
	require.NoError(s.T(), result.Error, "Failed to create tests attendance")
	return attendance
}

// Complete Workflow Test
func (s *TestSuite) TestCompletePayrollWorkflow() {
	// 1. Admin creates payroll period
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	periodData := map[string]string{
		"start_date": startDate,
		"end_date":   endDate,
	}

	t := s.T()
	w := s.makeRequest("POST", "/api/admin/payroll-periods", periodData, s.adminToken)
	if w.Code != http.StatusCreated {
		t.Logf("Payload being sent: %+v", periodData)
		t.Logf("token: %s", s.adminToken)
		t.Logf("Unexpected response status: %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}
	require.Equal(s.T(), http.StatusCreated, w.Code, "Failed to create payroll period")

	var periodResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &periodResp)
	require.NoError(s.T(), err, "Failed to parse period response")
	dataMap, _ := periodResp["data"].(map[string]interface{})
	periodID := uint(dataMap["id"].(float64))

	//2. Employee submits attendance
	date := time.Now().AddDate(0, -1, 5)
	for date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		date = date.AddDate(0, 0, -1) // mundurkan satu hari
	}
	attendanceDate := date.Format("2006-01-02")
	attendanceData := map[string]string{
		"date":      attendanceDate,
		"check_in":  "09:00:00",
		"check_out": "17:00:00",
	}

	w = s.makeRequest("POST", "/api/employee/attendance", attendanceData, s.employeeToken)
	if w.Code != http.StatusCreated {
		t.Logf("Payload being sent: %+v", attendanceData)
		t.Logf("token: %s", s.employeeToken)
		t.Logf("Unexpected response status: %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}
	assert.Equal(s.T(), http.StatusCreated, w.Code, "Failed to create attendance")

	// 3. Employee submits overtime
	overtimeData := map[string]interface{}{
		"date":        attendanceDate,
		"hours":       2.5,
		"description": "Project deadline",
	}

	w = s.makeRequest("POST", "/api/employee/overtime", overtimeData, s.employeeToken)
	assert.Equal(s.T(), http.StatusCreated, w.Code, "Failed to create overtime")

	// 4. Admin runs payroll
	runData := map[string]interface{}{
		"payroll_period_id": periodID,
	}

	w = s.makeRequest("POST", "/api/admin/payroll/run", runData, s.adminToken)
	assert.Equal(s.T(), http.StatusOK, w.Code, "Failed to run payroll")

	// 5. Employee views payslip
	w = s.makeRequest("GET", fmt.Sprintf("/api/employee/payslip?period_id=%d", periodID), nil, s.employeeToken)
	assert.Equal(s.T(), http.StatusOK, w.Code, "Failed to get payslip")

	var payslip map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &payslip)
	assert.NoError(s.T(), err, "Failed to parse payslip response")
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(TestSuite))
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}
