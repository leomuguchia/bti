package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"banking/models"
)

// These three helpers duplicate a slice of the repo's SQL because they need to
// run against a *sql.Tx rather than the repo's *sql.DB. Threading a Tx
func getAccountTx(tx *sql.Tx, id string) (*models.Account, error) {
	var a models.Account
	err := tx.QueryRow(
		`SELECT id, owner, balance, status, created_at, updated_at, last_txn_at
		 FROM accounts WHERE id = ?`, id,
	).Scan(&a.ID, &a.Owner, &a.Balance, &a.Status, &a.CreatedAt, &a.UpdatedAt, &a.LastTxnAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func saveAccountTx(tx *sql.Tx, a *models.Account) error {
	res, err := tx.Exec(
		`UPDATE accounts SET balance=?, updated_at=?, last_txn_at=? WHERE id=?`,
		a.Balance, a.UpdatedAt, a.LastTxnAt, a.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func appendLedgerTx(tx *sql.Tx, accountID, operation string, amount int64, meta string) error {
	_, err := tx.Exec(
		`INSERT INTO ledger_entries (id, account_id, operation, amount, meta, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		ulid.Make().String(), accountID, operation, amount, meta, time.Now().UTC(),
	)
	return err
}
