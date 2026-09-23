package account

import (
	"database/sql"

	"banking/models"
	"banking/services/ledger"
	"banking/utils"
)

type Repository interface {
	Create(a *models.Account) error
	Get(id string) (*models.Account, error)
	Save(a *models.Account) error
	Delete(id string) error
}

type Service interface {
	CreateAccount(owner string) (*models.Account, error)
	GetAccount(id string) (*models.Account, error)
	DeleteDormantAccount(id string) error
	Deposit(accountID string, amount int64) error
	Withdraw(accountID string, amount int64) error
	Transfer(fromID, toID string, amount int64) error
	CreditLoan(accountID string, amount int64) error
}

type service struct {
	repo   Repository
	ledger ledger.Service
	locks  *utils.KeyedMutex
	db     *sql.DB
}

func NewAccountService(repo Repository, ledgerSvc ledger.Service, db *sql.DB) Service {
	return &service{
		repo:   repo,
		ledger: ledgerSvc,
		locks:  utils.NewKeyedMutex(),
		db:     db,
	}
}
