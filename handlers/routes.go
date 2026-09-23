package handlers

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler, apiKey gin.HandlerFunc, rateLimit gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), rateLimit)

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	authed := r.Group("/", apiKey)
	{
		authed.POST("/accounts", h.CreateAccount)
		authed.GET("/accounts/:id", h.GetAccount)
		authed.DELETE("/accounts/:id", h.DeleteDormantAccount)
		authed.GET("/accounts/:id/ledger", h.LedgerHistory)
		authed.POST("/accounts/:id/deposit", h.Deposit)
		authed.POST("/accounts/:id/withdraw", h.Withdraw)
		authed.POST("/transfer", h.Transfer)

		authed.POST("/loans", h.CreateLoan)
		authed.POST("/loans/:id/disburse", h.DisburseLoan)
	}
	return r
}
