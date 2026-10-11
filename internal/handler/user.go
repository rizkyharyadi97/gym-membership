package handler

import (
	"gym-membership/internal/integration"
	"gym-membership/internal/model"
	"gym-membership/internal/service"
	"gym-membership/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
	email   *integration.ResendClient
}

func NewUserHandler(s service.UserService, email *integration.ResendClient) *UserHandler {
	return &UserHandler{service: s, email: email}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type depositRequest struct {
	Amount float64 `json:"amount"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}
	u := &model.User{Username: req.Username, Email: req.Email, Password: req.Password}
	if err := h.service.Register(c.Request().Context(), u); err != nil {
		return handleError(c, err)
	}
	u.Password = ""
	if h.email != nil {
		if err := h.email.Send(c.Request().Context(), u.Email,
			"Selamat datang di Gym Membership", "<p>Akun Gym Membership Anda berhasil dibuat. Selamat berlatih!</p>"); err != nil {
			log.Printf("Resend registration email failed: %v", err)
		}
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "success", "message": "user registered successfully", "data": u,
	})
}
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": "invalid user id",
		})
	}
	u, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": u})
}
func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": users})
}
func (h *UserHandler) AddDeposit(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "error", "message": "user authentication is required",
		})
	}
	var req depositRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}
	if err := h.service.AddDeposit(c.Request().Context(), uid, req.Amount); err != nil {
		return handleError(c, err)
	}
	u, err := h.service.GetByID(c.Request().Context(), uid)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success", "message": "deposit added successfully", "data": u,
	})
}
