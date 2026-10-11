package handler

import (
	"gym-membership/internal/service"
	"gym-membership/middleware"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct{ service service.TransactionService }

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

type purchaseRequest struct {
	MembershipID uint `json:"membership_id"`
}

func (h *TransactionHandler) PurchaseMembership(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(401, map[string]interface{}{
			"status": "error", "message": "user authentication is required",
		})
	}

	var r purchaseRequest
	if e := c.Bind(&r); e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid request body",
		})
	}

	t, e := h.service.PurchaseMembership(c.Request().Context(), uid, r.MembershipID)
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(201, map[string]interface{}{
		"status": "success", "message": "membership purchased successfully", "data": t,
	})
}
func (h *TransactionHandler) GetByID(c echo.Context) error {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		return c.JSON(400, map[string]interface{}{
			"status": "error", "message": "invalid transaction id",
		})
	}

	t, e := h.service.GetByID(c.Request().Context(), uint(id))
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{
		"status": "success", "data": t,
	})
}
func (h *TransactionHandler) GetByUserID(c echo.Context) error {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(401, map[string]interface{}{
			"status": "error", "message": "user authentication is required",
		})
	}

	t, e := h.service.GetByUserID(c.Request().Context(), uid)
	if e != nil {
		return handleError(c, e)
	}

	return c.JSON(200, map[string]interface{}{"status": "success", "data": t})
}
