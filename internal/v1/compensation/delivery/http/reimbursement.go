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

type ReimbursementHandler struct {
	uc     usecase_interface.IReimbursementUsecase
	config config.Config
}

func NewReimbursementHandler(e *echo.Group, config config.Config, uc usecase_interface.IReimbursementUsecase) {
	handler := &ReimbursementHandler{
		uc:     uc,
		config: config,
	}

	v1Compensation := e.Group("/compensation")

	v1Compensation.Use(jwt.JWTMiddleware)

	v1Compensation.POST("/reimbursement", handler.Create)
}

type ReimbursementRequest struct {
	Date        string  `json:"date"` // Format: YYYY-MM-DD
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (h *ReimbursementHandler) Create(c echo.Context) error {
	var req ReimbursementRequest

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

	overtime := &dto.CreateReimbursementRequest{
		UserID:      userID,
		EmployeeID:  employeeID,
		Date:        parsedDate,
		Description: req.Description,
		Amount:      req.Amount,
	}

	if err := h.uc.CreateReimbursement(ctx, overtime); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "reimbursement created", "request_id": requestID})
}
