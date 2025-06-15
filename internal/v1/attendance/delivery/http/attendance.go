package delivery_http

import (
	"net/http"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface"
	"github.com/labstack/echo/v4"
)

type AttendanceHandler struct {
	uc     usecase_interface.IAttendanceUsecase
	config config.Config
}

func NewAttendanceHandler(e *echo.Group, config config.Config, uc usecase_interface.IAttendanceUsecase) {
	handler := &AttendanceHandler{
		uc:     uc,
		config: config,
	}

	v1Attendance := e.Group("/attendances")

	v1Attendance.Use(jwt.JWTMiddleware)

	v1Attendance.POST("", handler.Create)
}

type AttendanceRequest struct {
	EmployeeID int    `json:"employee_id"`
	Date       string `json:"date"` // Format: YYYY-MM-DD
}

func (h *AttendanceHandler) Create(c echo.Context) error {
	var req AttendanceRequest

	ctx := c.Request().Context()
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	employeeID, _ := c.Get(jwt.EMPLOYEE_ID_KEY).(int)
	userID, _ := c.Get(jwt.USER_ID_KEY).(int)

	attendance := &dto.CreateAttendanceRequest{
		EmployeeID: employeeID,
		Date:       parsedDate,
		UserID:     userID,
	}

	if err := h.uc.CreateAttendance(ctx, attendance); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "attendance created"})
}
