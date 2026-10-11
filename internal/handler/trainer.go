package handler

import (
	"gym-membership/internal/model"
	"gym-membership/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TrainerHandler struct {
	service service.TrainerService
}

func NewTrainerHandler(s service.TrainerService) *TrainerHandler {
	return &TrainerHandler{service: s}
}

type trainerRequest struct {
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

func (h *TrainerHandler) Create(c echo.Context) error {

	var r trainerRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	t := &model.Trainer{Name: r.Name, Specialization: r.Specialization}
	if err := h.service.Create(c.Request().Context(), t); err != nil {
		return handleError(c, err)
	}

	return c.JSON(201, map[string]interface{}{"status": "success", "data": t})
}

func (h *TrainerHandler) GetByID(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid trainer id",
		})
	}

	t, e := h.service.GetByID(c.Request().Context(), uint(id))
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": t})
}

func (h *TrainerHandler) GetAll(c echo.Context) error {
	t, e := h.service.GetAll(c.Request().Context())
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": t})
}

func (h *TrainerHandler) Update(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid trainer id",
		})
	}

	var r trainerRequest
	if e = c.Bind(&r); e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	t := &model.Trainer{ID: uint(id), Name: r.Name, Specialization: r.Specialization}
	if e = h.service.Update(c.Request().Context(), t); e != nil {
		return handleError(c, e)
	}
	return c.JSON(200, map[string]interface{}{"status": "success", "data": t})
}

func (h *TrainerHandler) Delete(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid trainer id",
		})
	}

	if e = h.service.Delete(c.Request().Context(), uint(id)); e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{
		"status": "success", "message": "trainer deleted successfully",
	})
}
