package handlers

import (
	"errors"
	"net/http"

	"banking/models"
	"banking/services/account"
	"banking/services/loan"
)

// statusFor maps domain errors to HTTP statuses once, so handlers don't each
// reimplement the same switch.
func statusFor(err error) int {
	switch {
	case errors.Is(err, models.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, account.ErrHasBalance),
		errors.Is(err, account.ErrNotDormant),
		errors.Is(err, loan.ErrAlreadyDisbursed):
		return http.StatusConflict
	case errors.Is(err, account.ErrInvalidAmount),
		errors.Is(err, account.ErrInsufficientFunds),
		errors.Is(err, account.ErrSameAccount),
		errors.Is(err, loan.ErrInvalidAmount):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
