package handler

import (
	"gym-membership/internal/auth"
	"gym-membership/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userService service.UserService
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(us service.UserService, jm *auth.JWTManager) *AuthHandler {
	return &AuthHandler{userService: us, jwtManager: jm}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	user, err := h.userService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return handleError(c, err)
	}

	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "error", "message": "failed to generate authentication token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success", "data": map[string]interface{}{
			"token": token, "user": map[string]interface{}{
				"id": user.ID, "username": user.Username, "email": user.Email,
			},
		},
	})
}
