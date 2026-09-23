package handlers

import (
	"banking/services/account"
	"banking/services/ledger"
	"banking/services/loan"
)

type Handler struct {
	account account.Service
	loan    loan.Service
	ledger  ledger.Service
}

func NewHandler(accountSvc account.Service, loanSvc loan.Service, ledgerSvc ledger.Service) *Handler {
	return &Handler{account: accountSvc, loan: loanSvc, ledger: ledgerSvc}
}
