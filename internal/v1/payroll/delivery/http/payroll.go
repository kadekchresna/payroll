package delivery_http

import (
	"net/http"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/helper/logger"

	"github.com/kadekchresna/payroll/internal/v1/payroll/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/payroll/usecase/interface"

	"github.com/labstack/echo/v4"
)

type PayrollHandler struct {
	uc     usecase_interface.IPayrollUsecase
	config config.Config
}

func NewPayrollHandler(e *echo.Group, config config.Config, uc usecase_interface.IPayrollUsecase) {
	handler := &PayrollHandler{
		uc:     uc,
		config: config,
	}

	v1Compensation := e.Group("/payroll")

	v1Compensation.Use(jwt.JWTMiddleware)

	v1Compensation.POST("", handler.Create)
	v1Compensation.GET("", handler.GetPayroll)
	v1Compensation.GET("/summary", handler.GetPayrollSummary)
}

type PayrollRequest struct {
	AttendancePeriodID int `json:"attendance_period_id"`
}

func (h *PayrollHandler) Create(c echo.Context) error {
	var req PayrollRequest

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	userID, _ := c.Get(jwt.USER_ID_KEY).(int)

	payroll := dto.CreatePayrollRequest{
		UserID:             userID,
		AttendancePeriodID: req.AttendancePeriodID,
	}

	userRole, _ := c.Get(jwt.USER_ROLE_KEY).(string)

	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "only admin is authorized to access this content", "request_id": requestID})

	}

	if err := h.uc.CreatePayroll(ctx, payroll); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "create payroll failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "payroll created", "request_id": requestID})
}

type GetPayrollRequest struct {
	PayslipID int `query:"payslip_id"`
}

func (h *PayrollHandler) GetPayroll(c echo.Context) error {
	var req GetPayrollRequest

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	employeeID, _ := c.Get(jwt.EMPLOYEE_ID_KEY).(int)

	payroll := dto.GetEmployeePayrollRequest{
		EmployeeID: employeeID,
		PayslipID:  req.PayslipID,
	}

	res, err := h.uc.GetPayrollByID(ctx, &payroll)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "retrieve payroll failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "payroll retrieve successfully", "request_id": requestID, "data": res})
}

func (h *PayrollHandler) GetPayrollSummary(c echo.Context) error {
	ctx := c.Request().Context()

	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	userRole, _ := c.Get(jwt.USER_ROLE_KEY).(string)

	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "only admin is authorized to access this content", "request_id": requestID})

	}

	res, err := h.uc.GetPayrollSummary(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "retrieve payroll summary failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "payroll summary retrieve successfully", "request_id": requestID, "data": res})
}
