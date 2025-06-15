package jwt

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	USER_ID_KEY           = "user_id"
	USER_ROLE_KEY         = "user_role"
	EMPLOYEE_ID_KEY       = "employee_id"
	EMPLOYEE_FULLNAME_KEY = "employee_fullname"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token format"})
		}

		claims, err := ParseAccessToken(os.Getenv("APP_JWT_SECRET"), parts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		c.Set(USER_ID_KEY, claims.UserID)
		c.Set(USER_ROLE_KEY, claims.UserRole)
		c.Set(EMPLOYEE_ID_KEY, claims.EmployeeID)
		c.Set(EMPLOYEE_FULLNAME_KEY, claims.EmployeeFullname)

		return next(c)
	}
}
