package delivery_http

import (
	"net/http"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/helper/logger"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface"

	"github.com/labstack/echo/v4"
)

type AttendancePeriodHandler struct {
	uc     usecase_interface.IAttendancePeriodUsecase
	config config.Config
}

func NewAttendancePeriodHandler(e *echo.Group, config config.Config, uc usecase_interface.IAttendancePeriodUsecase) {
	handler := &AttendancePeriodHandler{
		uc:     uc,
		config: config,
	}

	v1Attendance := e.Group("/attendances-period")

	v1Attendance.Use(jwt.JWTMiddleware)

	v1Attendance.POST("", handler.Create)
}

type AttendancePeriodRequest struct {
	PeriodStart string `json:"period_start"` // Format: YYYY-MM-DD
	PeriodEnd   string `json:"period_end"`   // Format: YYYY-MM-DD

}

func (h *AttendancePeriodHandler) Create(c echo.Context) error {
	var req AttendancePeriodRequest

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	parsedPeriodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid period start format, use YYYY-MM-DD", "request_id": requestID})
	}

	parsedPeriodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid period end format, use YYYY-MM-DD", "request_id": requestID})
	}

	userID, _ := c.Get(jwt.USER_ID_KEY).(int)
	userRole, _ := c.Get(jwt.USER_ROLE_KEY).(string)

	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "only admin is authorized to access this content", "request_id": requestID})

	}

	attendance := &dto.CreateAttendancePeriodRequest{
		UserID:      userID,
		PeriodStart: parsedPeriodStart,
		PeriodEnd:   parsedPeriodEnd,
	}

	if err := h.uc.CreateAttendancePeriod(ctx, attendance); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"message": "create attendance period failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "create attendance period successfull", "request_id": requestID})
}
