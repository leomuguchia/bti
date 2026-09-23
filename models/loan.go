package models

import "time"

type LoanStatus string

const (
	LoanStatusActive LoanStatus = "active"
	LoanStatusPaid   LoanStatus = "paid"
)

type Loan struct {
	ID          string     `json:"id"`
	AccountID   string     `json:"account_id"`
	Principal   int64      `json:"principal"`
	Outstanding int64      `json:"outstanding"`
	Status      LoanStatus `json:"status"`
	DisbursedAt time.Time  `json:"disbursed_at"`
}
