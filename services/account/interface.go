package account

import (
	"database/sql"

	"banking/database"
	"banking/models"
	"banking/services/ledger"
	"banking/utils"
)

type Service interface {
	CreateAccount(owner models.Owner) (*models.Account, error)
	Login(nationalID, phoneNumber string) (*models.Account, error)
	GetAccount(id string) (*models.Account, error)
	DeleteDormantAccount(id string) error
	Deposit(accountID string, amount int64) error
	Withdraw(accountID string, amount int64) error
	Transfer(fromID, toID string, amount int64) error
	CreditLoan(accountID string, amount int64) error
}

type service struct {
	repo   database.Repository
	ledger ledger.Service
	locks  *utils.KeyedMutex
	db     *sql.DB
}

func NewAccountService(repo database.Repository, ledgerSvc ledger.Service, db *sql.DB) Service {
	return &service{
		repo:   repo,
		ledger: ledgerSvc,
		locks:  utils.NewKeyedMutex(),
		db:     db,
	}
}
