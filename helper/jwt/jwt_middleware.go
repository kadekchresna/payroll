package jwt

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
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

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.UserRole)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("employee_fullname", claims.EmployeeFullname)

		return next(c)
	}
}
