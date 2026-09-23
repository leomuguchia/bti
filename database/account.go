package database

import (
	"database/sql"
	"errors"

	"banking/models"
)

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepo {
	return &accountRepo{db: db}
}

const accountCols = `id, owner, balance, status, created_at, updated_at, last_txn_at`

func scanAccount(row interface{ Scan(...any) error }) (*models.Account, error) {
	var a models.Account
	err := row.Scan(&a.ID, &a.Owner, &a.Balance, &a.Status,
		&a.CreatedAt, &a.UpdatedAt, &a.LastTxnAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepo) Create(a *models.Account) error {
	_, err := r.db.Exec(
		`INSERT INTO accounts (id, owner, balance, status, created_at, updated_at, last_txn_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.Owner, a.Balance, a.Status, a.CreatedAt, a.UpdatedAt, a.LastTxnAt,
	)
	return err
}

func (r *accountRepo) Get(id string) (*models.Account, error) {
	return scanAccount(r.db.QueryRow(
		`SELECT `+accountCols+` FROM accounts WHERE id = ?`, id))
}

func (r *accountRepo) Save(a *models.Account) error {
	res, err := r.db.Exec(
		`UPDATE accounts SET owner=?, balance=?, status=?, updated_at=?, last_txn_at=?
		 WHERE id=?`,
		a.Owner, a.Balance, a.Status, a.UpdatedAt, a.LastTxnAt, a.ID,
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

func (r *accountRepo) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM accounts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *accountRepo) List() ([]*models.Account, error) {
	rows, err := r.db.Query(`SELECT ` + accountCols + ` FROM accounts ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Account
	for rows.Next() {
		a, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// GetTx / SaveTx are the transaction-aware variants used by Transfer so the
// debit+credit+ledger write commit as one unit.
func (r *accountRepo) GetTx(tx *sql.Tx, id string) (*models.Account, error) {
	return scanAccount(tx.QueryRow(
		`SELECT `+accountCols+` FROM accounts WHERE id = ?`, id))
}

func (r *accountRepo) SaveTx(tx *sql.Tx, a *models.Account) error {
	res, err := tx.Exec(
		`UPDATE accounts SET owner=?, balance=?, status=?, updated_at=?, last_txn_at=?
		 WHERE id=?`,
		a.Owner, a.Balance, a.Status, a.UpdatedAt, a.LastTxnAt, a.ID,
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
