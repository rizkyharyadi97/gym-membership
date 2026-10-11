package middleware

import (
	"gym-membership/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const UserIDKey = "user_id"

func JWTMiddleware(manager *auth.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": "error", "message": "invalid authorization header",
				})
			}
			claims, err := manager.ValidateToken(strings.TrimSpace(parts[1]))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"status": "error", "message": "invalid or expired token",
				})
			}
			c.Set(UserIDKey, claims.UserID)
			return next(c)
		}
	}
}

func GetUserID(c echo.Context) (uint, bool) {
	id, ok := c.Get(UserIDKey).(uint)
	return id, ok && id > 0
}
