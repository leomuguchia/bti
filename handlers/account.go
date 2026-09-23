package handlers

import (
	"banking/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var req struct {
		Owner models.Owner `json:"owner" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "owner details are required")
		return
	}
	if err := req.Owner.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.account.CreateAccount(req.Owner)
	if err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	respondJSON(c, http.StatusCreated, a)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		NationalID  string `json:"national_id" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "national ID and phone number are required")
		return
	}
	a, err := h.account.Login(req.NationalID, req.PhoneNumber)
	if err != nil {
		// generic message, do not confirm whether the ID exists
		respondError(c, http.StatusUnauthorized, "no account matches those details")
		return
	}
	respondJSON(c, http.StatusOK, a)
}

func (h *Handler) GetAccount(c *gin.Context) {
	a, err := h.account.GetAccount(c.Param("id"))
	if err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	respondJSON(c, http.StatusOK, a)
}

func (h *Handler) DeleteDormantAccount(c *gin.Context) {
	if err := h.account.DeleteDormantAccount(c.Param("id")); err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Deposit(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.account.Deposit(c.Param("id"), req.Amount); err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Withdraw(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.account.Withdraw(c.Param("id"), req.Amount); err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) Transfer(c *gin.Context) {
	var req struct {
		FromID string `json:"from_id" binding:"required"`
		ToID   string `json:"to_id" binding:"required"`
		Amount int64  `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.account.Transfer(req.FromID, req.ToID, req.Amount); err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) LedgerHistory(c *gin.Context) {
	entries, err := h.ledger.History(c.Param("id"))
	if err != nil {
		respondError(c, statusFor(err), err.Error())
		return
	}
	respondJSON(c, http.StatusOK, entries)
}
