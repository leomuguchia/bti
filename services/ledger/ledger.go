package ledger

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"

	"banking/models"
)

func newEntry(accountID, operation string, amount int64, meta string) *models.LedgerEntry {
	return &models.LedgerEntry{
		ID:        ulid.Make().String(),
		AccountID: accountID,
		Operation: operation,
		Amount:    amount,
		Meta:      meta,
		CreatedAt: time.Now().UTC(),
	}
}

func (s *service) Record(accountID, operation string, amount int64, meta string) error {
	return s.repo.Append(newEntry(accountID, operation, amount, meta))
}

func (s *service) RecordTx(tx *sql.Tx, accountID, operation string, amount int64, meta string) error {
	return s.repo.AppendTx(tx, newEntry(accountID, operation, amount, meta))
}

func (s *service) History(accountID string) ([]*models.LedgerEntry, error) {
	return s.repo.ListByAccount(accountID)
}
