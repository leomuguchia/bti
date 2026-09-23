package database

import (
	"database/sql"

	"banking/models"
)

type ledgerRepo struct {
	db *sql.DB
}

func NewLedgerRepository(db *sql.DB) *ledgerRepo {
	return &ledgerRepo{db: db}
}

func (r *ledgerRepo) Append(e *models.LedgerEntry) error {
	_, err := r.db.Exec(
		`INSERT INTO ledger_entries (id, account_id, operation, amount, meta, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		e.ID, e.AccountID, e.Operation, e.Amount, e.Meta, e.CreatedAt,
	)
	return err
}

// AppendTx is the transaction-aware variant used by Transfer.
func (r *ledgerRepo) AppendTx(tx *sql.Tx, e *models.LedgerEntry) error {
	_, err := tx.Exec(
		`INSERT INTO ledger_entries (id, account_id, operation, amount, meta, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		e.ID, e.AccountID, e.Operation, e.Amount, e.Meta, e.CreatedAt,
	)
	return err
}

func (r *ledgerRepo) ListByAccount(accountID string) ([]*models.LedgerEntry, error) {
	rows, err := r.db.Query(
		`SELECT id, account_id, operation, amount, meta, created_at
		 FROM ledger_entries WHERE account_id = ? ORDER BY created_at`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.LedgerEntry
	for rows.Next() {
		var e models.LedgerEntry
		if err := rows.Scan(&e.ID, &e.AccountID, &e.Operation, &e.Amount, &e.Meta, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}
