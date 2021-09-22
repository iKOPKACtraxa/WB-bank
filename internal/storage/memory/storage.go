package memorystorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/iKOPKACtraxa/wb-bank/internal/storage"
)

var (
	ErrNotEnoughBalance = errors.New("account has not anough balance")
	ErrAccIsAbsent      = errors.New("account is not exist")
	ErrAccNameInvalid   = errors.New("account name is not valid")
)

type StorageInMemory struct {
	m sync.Map
}

// New returns a StorageInMemory.
func New() *StorageInMemory {
	return &StorageInMemory{
		m: sync.Map{},
	}
}

// Set sets a new account with balance.
func (s *StorageInMemory) Set(account storage.AccountID, balance storage.AccountBalance) error {
	if account == "" {
		return ErrAccNameInvalid
	}
	s.m.Store(account, balance)
	return nil
}

// Get reads a balance of account.
func (s *StorageInMemory) Get(account storage.AccountID) (storage.AccountBalance, error) {
	balance, ok := s.m.Load(account)
	if !ok {
		return 0, ErrAccIsAbsent
	}
	return balance.(storage.AccountBalance), nil
}

// Transfer ransfers money from accountFrom to accountTo using the amount.
func (s *StorageInMemory) Transfer(accountFrom, accountTo storage.AccountID, amountToTransfer storage.AccountBalance) error {
	balanceForAccountFrom, err := s.Get(accountFrom)
	if err != nil {
		return err
	}
	balanceForAccountTo, err := s.Get(accountTo)
	if err != nil {
		return err
	}
	if balanceForAccountFrom < amountToTransfer {
		return ErrNotEnoughBalance
	}
	err = s.Set(accountFrom, balanceForAccountFrom-amountToTransfer)
	if err != nil {
		return err
	}
	err = s.Set(accountTo, balanceForAccountTo+amountToTransfer)
	if err != nil {
		return err
	}
	return nil
}

// Print метод только для отладки todo удалить реализацию и из интерфейса storage.Storage.
func (s *StorageInMemory) Print() {
	mapToPrint := map[string]interface{}{}
	s.m.Range(func(key, value interface{}) bool {
		mapToPrint[fmt.Sprint(key)] = value
		return true
	})
	fmt.Println("Текущее состояние счетов: ", mapToPrint)
}
