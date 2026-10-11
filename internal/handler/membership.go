package handler

import (
	"gym-membership/internal/model"
	"gym-membership/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MembershipHandler struct{ service service.MembershipService }

func NewMembershipHandler(s service.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: s}
}

type membershipRequest struct {
	Name         string  `json:"name"`
	Costs        float64 `json:"costs"`
	DurationDays int     `json:"duration_days"`
	Availability bool    `json:"availability"`
	Category     string  `json:"category"`
}

func (h *MembershipHandler) Create(c echo.Context) error {

	var r membershipRequest

	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	m := &model.Membership{Name: r.Name, Costs: r.Costs, DurationDays: r.DurationDays, Availability: r.Availability, Category: r.Category}
	if err := h.service.Create(c.Request().Context(), m); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"status": "success", "data": m})
}

func (h *MembershipHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid membership id",
		})
	}

	m, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": m})
}

func (h *MembershipHandler) GetAll(c echo.Context) error {
	m, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": m})
}

func (h *MembershipHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid membership id",
		})
	}

	var r membershipRequest
	if err = c.Bind(&r); err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	m := &model.Membership{ID: uint(id), Name: r.Name, Costs: r.Costs, DurationDays: r.DurationDays, Availability: r.Availability, Category: r.Category}
	if err = h.service.Update(c.Request().Context(), m); err != nil {
		return handleError(c, err)
	}
	return c.JSON(200, map[string]interface{}{"status": "success", "data": m})
}

func (h *MembershipHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid membership id",
		})
	}

	if err = h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return handleError(c, err)
	}

	return c.JSON(200, map[string]interface{}{
		"status": "success", "message": "membership deleted successfully",
	})
}
