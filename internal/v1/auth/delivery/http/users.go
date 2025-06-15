package delivery_http

import (
	"net/http"

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
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	if err := h.uc.Create(c.Request().Context(), req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Register successful"})
}

func (h *UsersHandler) Login(c echo.Context) error {
	var req dto.LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	}

	res, err := h.uc.Login(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "Login Successful", "data": res})
}
