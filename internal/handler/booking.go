package handler

import (
	"gym-membership/internal/service"
	"gym-membership/middleware"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(s service.BookingService) *BookingHandler {
	return &BookingHandler{service: s}
}

type createBookingRequest struct {
	GymClassID uint `json:"gym_class_id"`
}

func (h *BookingHandler) Create(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(401, map[string]interface{}{
			"status": "error", "message": "user authentication is required",
		})
	}

	var r createBookingRequest
	if e := c.Bind(&r); e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	b, e := h.service.Create(c.Request().Context(), uid, r.GymClassID)
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(201, map[string]interface{}{"status": "success", "data": b})
}

func (h *BookingHandler) GetByID(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid booking id",
		})
	}

	b, e := h.service.GetByID(c.Request().Context(), uint(id))
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": b})
}

func (h *BookingHandler) GetByUserID(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(401, map[string]interface{}{
			"status": "error", "message": "user authentication is required",
		})
	}

	b, e := h.service.GetByUserID(c.Request().Context(), uid)
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": b})
}

func (h *BookingHandler) GetByClassID(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("classID"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid class id",
		})
	}

	b, e := h.service.GetByClassID(c.Request().Context(), uint(id))
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": b})
}

func (h *BookingHandler) Cancel(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid booking id",
		})
	}

	if e = h.service.Cancel(c.Request().Context(), uint(id)); e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{
		"status": "success", "message": "booking cancelled successfully",
	})
}
