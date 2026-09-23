package database

import (
	"database/sql"
	"errors"

	"banking/models"
)

type loanRepo struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *loanRepo {
	return &loanRepo{db: db}
}

const loanCols = `id, account_id, principal, outstanding, status, disbursed_at`

func scanLoan(row interface{ Scan(...any) error }) (*models.Loan, error) {
	var l models.Loan
	var disbursed sql.NullTime
	err := row.Scan(&l.ID, &l.AccountID, &l.Principal, &l.Outstanding, &l.Status, &disbursed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if disbursed.Valid {
		l.DisbursedAt = disbursed.Time
	}
	return &l, nil
}

func (r *loanRepo) Create(l *models.Loan) error {
	var disbursed any
	if !l.DisbursedAt.IsZero() {
		disbursed = l.DisbursedAt
	}
	_, err := r.db.Exec(
		`INSERT INTO loans (id, account_id, principal, outstanding, status, disbursed_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		l.ID, l.AccountID, l.Principal, l.Outstanding, l.Status, disbursed,
	)
	return err
}

func (r *loanRepo) Get(id string) (*models.Loan, error) {
	return scanLoan(r.db.QueryRow(`SELECT `+loanCols+` FROM loans WHERE id = ?`, id))
}

func (r *loanRepo) Save(l *models.Loan) error {
	res, err := r.db.Exec(
		`UPDATE loans SET principal=?, outstanding=?, status=?, disbursed_at=? WHERE id=?`,
		l.Principal, l.Outstanding, l.Status, l.DisbursedAt, l.ID,
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

func (r *loanRepo) GetByAccount(accountID string) (*models.Loan, error) {
	return scanLoan(r.db.QueryRow(`SELECT `+loanCols+` FROM loans WHERE account_id = ?`, accountID))
}
