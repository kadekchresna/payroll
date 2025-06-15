package delivery_http

import (
	"net/http"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/helper/logger"
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/compensation/usecase/interface"

	"github.com/labstack/echo/v4"
)

type OvertimeHandler struct {
	uc     usecase_interface.IOvertimeUsecase
	config config.Config
}

func NewOvertimeHandler(e *echo.Group, config config.Config, uc usecase_interface.IOvertimeUsecase) {
	handler := &OvertimeHandler{
		uc:     uc,
		config: config,
	}

	v1Compensation := e.Group("/compensation")

	v1Compensation.Use(jwt.JWTMiddleware)

	v1Compensation.POST("/overtime", handler.Create)
}

type OvertimeRequest struct {
	Date  string `json:"date"` // Format: YYYY-MM-DD
	Hours int    `json:"hours"`
}

func (h *OvertimeHandler) Create(c echo.Context) error {
	var req OvertimeRequest

	ctx := c.Request().Context()
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	userID, _ := c.Get(jwt.USER_ID_KEY).(int)
	employeeID, _ := c.Get(jwt.EMPLOYEE_ID_KEY).(int)

	overtime := &dto.CreateOvertimeRequest{
		UserID:     userID,
		EmployeeID: employeeID,
		Date:       parsedDate,
		Hours:      req.Hours,
	}

	if err := h.uc.CreateOvertime(ctx, overtime); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "overtime created", "request_id": requestID})
}
