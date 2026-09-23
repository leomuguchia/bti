package account

import (
	"errors"
	"time"
)

var (
	ErrInvalidAmount     = errors.New("account: amount must be positive")
	ErrInsufficientFunds = errors.New("account: insufficient funds")
	ErrSameAccount       = errors.New("account: cannot transfer to the same account")
)

func (s *service) Deposit(accountID string, amount int64) error {
	return s.credit(accountID, amount, "deposit")
}

func (s *service) CreditLoan(accountID string, amount int64) error {
	return s.credit(accountID, amount, "loan_disbursement")
}

func (s *service) credit(accountID string, amount int64, op string) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	unlock := s.locks.Lock(accountID)
	defer unlock()

	a, err := s.repo.Get(accountID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	a.Balance += amount
	a.UpdatedAt = now
	a.LastTxnAt = now
	if err := s.repo.Save(a); err != nil {
		return err
	}
	return s.ledger.Record(accountID, op, amount, "")
}

func (s *service) Withdraw(accountID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	unlock := s.locks.Lock(accountID)
	defer unlock()

	a, err := s.repo.Get(accountID)
	if err != nil {
		return err
	}
	if a.Balance < amount {
		return ErrInsufficientFunds
	}
	now := time.Now().UTC()
	a.Balance -= amount
	a.UpdatedAt = now
	a.LastTxnAt = now
	if err := s.repo.Save(a); err != nil {
		return err
	}
	return s.ledger.Record(accountID, "withdraw", amount, "")
}

// Transfer locks both accounts in sorted order (no deadlock between opposite-
// direction concurrent transfers) and writes both balance updates plus both
// ledger entries in one SQL transaction so a partial transfer can't commit.
func (s *service) Transfer(fromID, toID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if fromID == toID {
		return ErrSameAccount
	}

	first, second := fromID, toID
	if second < first {
		first, second = second, first
	}
	unlock1 := s.locks.Lock(first)
	defer unlock1()
	unlock2 := s.locks.Lock(second)
	defer unlock2()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	from, err := getAccountTx(tx, fromID)
	if err != nil {
		return err
	}
	to, err := getAccountTx(tx, toID)
	if err != nil {
		return err
	}
	if from.Balance < amount {
		return ErrInsufficientFunds
	}

	now := time.Now().UTC()
	from.Balance -= amount
	from.UpdatedAt = now
	from.LastTxnAt = now
	to.Balance += amount
	to.UpdatedAt = now
	to.LastTxnAt = now

	if err := saveAccountTx(tx, from); err != nil {
		return err
	}
	if err := saveAccountTx(tx, to); err != nil {
		return err
	}
	if err := appendLedgerTx(tx, fromID, "transfer_out", amount, toID); err != nil {
		return err
	}
	if err := appendLedgerTx(tx, toID, "transfer_in", amount, fromID); err != nil {
		return err
	}
	return tx.Commit()
}
