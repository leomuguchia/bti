package loan

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"banking/models"
)

var (
	ErrAlreadyDisbursed  = errors.New("loan: already disbursed")
	ErrInvalidAmount     = errors.New("loan: amount must be positive")
	ErrLoanAlreadyExists = errors.New("loan: account already has a loan")
)

func (s *service) CreateLoanAccount(accountID string) (*models.Loan, error) {
	_, err := s.repo.GetByAccount(accountID)
	if err == nil {
		return nil, ErrLoanAlreadyExists
	}
	if !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}

	l := &models.Loan{
		ID:        ulid.Make().String(),
		AccountID: accountID,
		Status:    models.LoanStatusActive,
	}
	if err := s.repo.Create(l); err != nil {
		return nil, err
	}
	return l, nil
}

// Disburse credits the linked account via the account service — so the credit
// takes the same per-account lock and lands in the same ledger, tagged as a
// loan disbursement rather than a plain deposit. Known boundary: the credit
// and the loan-record save commit separately; a crash between them would
// leave the account credited and the loan un-updated.
func (s *service) Disburse(loanID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	l, err := s.repo.Get(loanID)
	if err != nil {
		return err
	}
	if !l.DisbursedAt.IsZero() {
		return ErrAlreadyDisbursed
	}
	if err := s.account.CreditLoan(l.AccountID, amount); err != nil {
		return err
	}
	l.Principal = amount
	l.Outstanding = amount
	l.DisbursedAt = time.Now().UTC()
	return s.repo.Save(l)
}
