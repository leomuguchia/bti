package ledger

import (
	"database/sql"

	"banking/models"
)

type Repository interface {
	Append(e *models.LedgerEntry) error
	AppendTx(tx *sql.Tx, e *models.LedgerEntry) error
	ListByAccount(accountID string) ([]*models.LedgerEntry, error)
}

type Service interface {
	Record(accountID, operation string, amount int64, meta string) error
	RecordTx(tx *sql.Tx, accountID, operation string, amount int64, meta string) error
	History(accountID string) ([]*models.LedgerEntry, error)
}

type service struct {
	repo Repository
}

func NewLedgerService(repo Repository) Service {
	return &service{repo: repo}
}
