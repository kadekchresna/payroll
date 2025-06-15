package delivery_http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/internal/v1/attendance/model"
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
	v1Attendance.GET("/:id", handler.Get)
	v1Attendance.PUT("/:id", handler.Update)
	v1Attendance.DELETE("/:id", handler.Delete)
}

type AttendanceRequest struct {
	EmployeeID int    `json:"employee_id"`
	Date       string `json:"date"` // Format: YYYY-MM-DD
	CreatedBy  int    `json:"created_by"`
	UpdatedBy  int    `json:"updated_by"`
}

func (h *AttendanceHandler) Create(c echo.Context) error {
	var req AttendanceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	attendance := &model.Attendance{
		EmployeeID: req.EmployeeID,
		Date:       parsedDate,
		CreatedBy:  req.CreatedBy,
		UpdatedBy:  req.UpdatedBy,
	}

	if err := h.uc.CreateAttendance(c.Request().Context(), attendance); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "attendance created"})
}

func (h *AttendanceHandler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	attendance, err := h.uc.GetAttendance(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if attendance == nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, attendance)
}

func (h *AttendanceHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var req AttendanceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
	}

	attendance := &model.Attendance{
		ID:         id,
		EmployeeID: req.EmployeeID,
		Date:       parsedDate,
		UpdatedBy:  req.UpdatedBy,
	}

	if err := h.uc.UpdateAttendance(c.Request().Context(), attendance); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "attendance updated"})
}

func (h *AttendanceHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.uc.DeleteAttendance(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "attendance deleted"})
}
