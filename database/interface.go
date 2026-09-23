package database

import (
	"database/sql"

	"banking/models"
)

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepo {
	return &accountRepo{db: db}
}

const accountCols = `id, owner, balance, status, created_at, updated_at, last_txn_at`

type Repository interface {
	Create(a *models.Account) error
	Get(id string) (*models.Account, error)
	Save(a *models.Account) error
	Delete(id string) error
	FindByIdentity(nationalID, phoneNumber string) (*models.Account, error)
}
