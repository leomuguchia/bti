package models

import "time"

type LedgerEntry struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Operation string    `json:"operation"` // deposit, withdraw, transfer_out, transfer_in, loan_disbursement
	Amount    int64     `json:"amount"`
	Meta      string    `json:"meta,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
