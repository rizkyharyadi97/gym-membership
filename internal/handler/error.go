package handler

import (
	"errors"
	"gym-membership/internal/repository"
	"gym-membership/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, repository.ErrUserNotFound),
		errors.Is(err, repository.ErrMembershipNotFound),
		errors.Is(err, repository.ErrTrainerNotFound),
		errors.Is(err, repository.ErrGymClassNotFound),
		errors.Is(err, repository.ErrTransactionNotFound),
		errors.Is(err, repository.ErrBookingNotFound):

		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	case errors.Is(err, service.ErrInvalidCredentials):

		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	case errors.Is(err, service.ErrUsernameRequired),
		errors.Is(err, service.ErrEmailRequired),
		errors.Is(err, service.ErrPasswordRequired),
		errors.Is(err, service.ErrMembershipNameRequired),
		errors.Is(err, service.ErrInvalidMembershipCost),
		errors.Is(err, service.ErrInvalidDuration),
		errors.Is(err, service.ErrTrainerNameRequired),
		errors.Is(err, service.ErrTrainerSpecializationRequired),
		errors.Is(err, service.ErrClassNameRequired),
		errors.Is(err, service.ErrInvalidQuota),
		errors.Is(err, service.ErrInvalidSchedule),
		errors.Is(err, service.ErrInvalidDeposit):

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	case errors.Is(err, service.ErrUserExists):

		return c.JSON(http.StatusConflict, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	case errors.Is(err, service.ErrInsufficientDeposit),
		errors.Is(err, service.ErrMembershipUnavailable):

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	case errors.Is(err, service.ErrClassUnavailable),
		errors.Is(err, service.ErrClassFull),
		errors.Is(err, service.ErrAlreadyBooked):

		return c.JSON(http.StatusConflict, map[string]interface{}{
			"status": "error", "message": err.Error(),
		})

	default:
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": "error", "message": "internal server error",
		})
	}
}
