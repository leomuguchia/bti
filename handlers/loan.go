package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateLoan(c *gin.Context) {
	var req struct {
		AccountID string `json:"account_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "account_id is required")
		return
	}
	l, err := h.loan.CreateLoanAccount(req.AccountID)
	if err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	respondJSON(c, http.StatusCreated, l)
}

func (h *Handler) DisburseLoan(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.loan.Disburse(c.Param("id"), req.Amount); err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	c.Status(http.StatusOK)
}
