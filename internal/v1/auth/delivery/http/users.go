package delivery_http

import (
	"net/http"

	"github.com/kadekchresna/payroll/helper/logger"
	"github.com/kadekchresna/payroll/internal/v1/auth/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/auth/usecase/interface"
	"github.com/labstack/echo/v4"
)

type UsersHandler struct {
	uc usecase_interface.IUserUsecase
}

func NewUsersHandler(e *echo.Group, uc usecase_interface.IUserUsecase) {
	handler := &UsersHandler{uc: uc}

	e.POST("/auth/register", handler.Register)
	e.POST("/auth/login", handler.Login)
}

func (h *UsersHandler) Register(c echo.Context) error {
	var req dto.CreateUserRequest

	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	if err := h.uc.Create(ctx, req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "register failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "register successful", "request_id": requestID})
}

func (h *UsersHandler) Login(c echo.Context) error {
	var req dto.LoginUserRequest
	ctx := c.Request().Context()
	requestID, _ := ctx.Value(logger.RequestIDKey).(string)
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "invalid input", "request_id": requestID})
	}

	res, err := h.uc.Login(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"message": "login failed", "error": err.Error(), "request_id": requestID})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "login successful", "data": res, "request_id": requestID})
}
