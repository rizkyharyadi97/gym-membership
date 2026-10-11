package handler

import (
	"gym-membership/internal/model"
	"gym-membership/internal/service"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type GymClassHandler struct {
	service service.GymClassService
}

func NewGymClassHandler(s service.GymClassService) *GymClassHandler {
	return &GymClassHandler{service: s}
}

type gymClassRequest struct {
	TrainerID    uint      `json:"trainer_id"`
	ClassName    string    `json:"class_name"`
	Schedule     time.Time `json:"schedule"`
	Quota        int       `json:"quota"`
	Availability bool      `json:"availability"`
}

func (h *GymClassHandler) Create(c echo.Context) error {
	var r gymClassRequest
	if e := c.Bind(&r); e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	g := &model.GymClass{TrainerID: r.TrainerID, ClassName: r.ClassName, Schedule: r.Schedule, Quota: r.Quota, Availability: r.Availability}
	if e := h.service.Create(c.Request().Context(), g); e != nil {
		return handleError(c, e)
	}

	return c.JSON(201, map[string]interface{}{"status": "success", "data": g})
}

func (h *GymClassHandler) GetByID(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid class id",
		})
	}

	g, e := h.service.GetByID(c.Request().Context(), uint(id))
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": g})
}

func (h *GymClassHandler) GetAll(c echo.Context) error {
	g, e := h.service.GetAll(c.Request().Context())
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": g})
}

func (h *GymClassHandler) Update(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid class id",
		})
	}

	var r gymClassRequest
	if e = c.Bind(&r); e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	g := &model.GymClass{ID: uint(id), TrainerID: r.TrainerID, ClassName: r.ClassName, Schedule: r.Schedule, Quota: r.Quota, Availability: r.Availability}
	if e = h.service.Update(c.Request().Context(), g); e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": g})
}

func (h *GymClassHandler) Delete(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid class id",
		})
	}

	if e = h.service.Delete(c.Request().Context(), uint(id)); e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{
		"status": "success", "message": "gym class deleted successfully",
	})
}
