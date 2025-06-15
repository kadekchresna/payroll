package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/logger"
	helper_time "github.com/kadekchresna/payroll/helper/time"
	driver_db "github.com/kadekchresna/payroll/infrastructure/db"
	helper_db "github.com/kadekchresna/payroll/infrastructure/db/helper"
	audit_repo "github.com/kadekchresna/payroll/internal/v1/audit/repository"
	employee_repo "github.com/kadekchresna/payroll/internal/v1/employee/repository"

	attendance_delivery_http "github.com/kadekchresna/payroll/internal/v1/attendance/delivery/http"
	attendance_repo "github.com/kadekchresna/payroll/internal/v1/attendance/repository"
	attendance_usecase "github.com/kadekchresna/payroll/internal/v1/attendance/usecase"

	auth_delivery_http "github.com/kadekchresna/payroll/internal/v1/auth/delivery/http"
	auth_repository "github.com/kadekchresna/payroll/internal/v1/auth/repository"
	auth_usecase "github.com/kadekchresna/payroll/internal/v1/auth/usecase"

	compensation_delivery_http "github.com/kadekchresna/payroll/internal/v1/compensation/delivery/http"
	compensation_repository "github.com/kadekchresna/payroll/internal/v1/compensation/repository"
	compensation_usecase "github.com/kadekchresna/payroll/internal/v1/compensation/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

const (
	STAGING     = `stg`
	PRODUCTIOON = `prd`
)

func init() {
	if os.Getenv("APP_ENV") != PRODUCTIOON {

		// init invoke env before everything
		cobra.OnInitialize(initConfig)

	}

	// adding command invokable
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "web",
	Short: "Running Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func initConfig() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("error load ENV, %s", err.Error()))
	}
}

func run() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(logger.RequestIDMiddleware)
	e.Use(logger.ClientIPMiddleware)

	config := config.InitConfig()
	db := driver_db.InitDB(config.DatabaseDSN)

	// V1 Endpoints
	v1 := e.Group("/api/v1")

	timer := helper_time.NewTime(nil)
	transactionBundler := helper_db.NewTransactionBundler(db)

	employeeRepo := employee_repo.NewEmployeeRepository(db)
	auditRepo := audit_repo.NewAuditRepository(db)

	attendanceRepo := attendance_repo.NewAttendanceRepository(db)
	attendanceUsecase := attendance_usecase.NewAttendanceUsecase(timer, attendanceRepo, employeeRepo, auditRepo, transactionBundler)
	attendance_delivery_http.NewAttendanceHandler(v1, config, attendanceUsecase)

	attendancePeriodRepo := attendance_repo.NewAttendancePeriodRepository(db)
	attendancePeriodUsecase := attendance_usecase.NewAttendancePeriodUsecase(timer, attendancePeriodRepo, transactionBundler, auditRepo)
	attendance_delivery_http.NewAttendancePeriodHandler(v1, config, attendancePeriodUsecase)

	userRepo := auth_repository.NewUserRepo(db)
	userUsecase := auth_usecase.NewUserUsecase(config, timer, userRepo, employeeRepo)
	auth_delivery_http.NewUsersHandler(v1, userUsecase)

	overtimeRepo := compensation_repository.NewOvertimeRepository(db)
	reimbursementRepo := compensation_repository.NewReimbursementRepository(db)
	overtimeUsecase := compensation_usecase.NewOvertimeUsecase(overtimeRepo, attendanceRepo, auditRepo, transactionBundler, timer)
	reimbursementUsecase := compensation_usecase.NewReimbursementUsecase(timer, reimbursementRepo, auditRepo, transactionBundler)
	compensation_delivery_http.NewOvertimeHandler(v1, config, overtimeUsecase)
	compensation_delivery_http.NewReimbursementHandler(v1, config, reimbursementUsecase)
	// V1 Endpoints

	s := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort),
		Handler: e,
	}

	logger.Log().Info(fmt.Sprintf("%s service started...", config.AppName))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	logger.Log().Info(fmt.Sprintf("%s service finished", config.AppName))
}
