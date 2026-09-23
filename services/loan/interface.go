package loan

import (
	"banking/models"
	"banking/services/account"
)

type Repository interface {
	Create(l *models.Loan) error
	Get(id string) (*models.Loan, error)
	Save(l *models.Loan) error
	GetByAccount(accountID string) (*models.Loan, error)
}

type Service interface {
	CreateLoanAccount(accountID string) (*models.Loan, error)
	Disburse(loanID string, amount int64) error
}

type service struct {
	repo    Repository
	account account.Service
}

func NewLoanService(repo Repository, accountSvc account.Service) Service {
	return &service{repo: repo, account: accountSvc}
}
