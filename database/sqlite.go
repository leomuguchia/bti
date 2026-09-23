package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string) (*sql.DB, error) {
	// WAL + busy_timeout so concurrent writers wait instead of erroring.
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id          TEXT PRIMARY KEY,
			owner       TEXT NOT NULL,
			balance     INTEGER NOT NULL DEFAULT 0,
			status      TEXT NOT NULL,
			created_at  TIMESTAMP NOT NULL,
			updated_at  TIMESTAMP NOT NULL,
			last_txn_at TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ledger_entries (
			id         TEXT PRIMARY KEY,
			account_id TEXT NOT NULL,
			operation  TEXT NOT NULL,
			amount     INTEGER NOT NULL,
			meta       TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_account_created
			ON ledger_entries(account_id, created_at);`,
		`CREATE TABLE IF NOT EXISTS loans (
			id           TEXT PRIMARY KEY,
			account_id   TEXT NOT NULL,
			principal    INTEGER NOT NULL DEFAULT 0,
			outstanding  INTEGER NOT NULL DEFAULT 0,
			status       TEXT NOT NULL,
			disbursed_at TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_loans_account
			ON loans(account_id);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
