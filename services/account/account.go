package account

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"

	"banking/models"
)

const DormancyPeriod = 90 * 24 * time.Hour

var (
	ErrNotDormant = errors.New("account: not eligible for deletion, not dormant")
	ErrHasBalance = errors.New("account: cannot delete account with non-zero balance")
)

func (s *service) CreateAccount(owner string) (*models.Account, error) {
	now := time.Now().UTC()
	a := &models.Account{
		ID:        ulid.Make().String(),
		Owner:     owner,
		Balance:   0,
		Status:    models.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		LastTxnAt: now,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *service) GetAccount(id string) (*models.Account, error) {
	return s.repo.Get(id)
}

// DeleteDormantAccount holds the account lock across check-and-delete so a
// concurrent deposit can't land between the balance check and the delete.
func (s *service) DeleteDormantAccount(id string) error {
	unlock := s.locks.Lock(id)
	defer unlock()

	a, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if a.Balance != 0 {
		return ErrHasBalance
	}
	if time.Since(a.LastTxnAt) < DormancyPeriod {
		return ErrNotDormant
	}
	return s.repo.Delete(id)
}
