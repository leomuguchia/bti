package models

import "time"

type Status string

const (
	StatusActive  Status = "active"
	StatusDormant Status = "dormant"
	StatusClosed  Status = "closed"
)

type Account struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"` // minor units
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastTxnAt time.Time `json:"last_txn_at"`
}
