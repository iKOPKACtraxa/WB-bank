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

type balance struct {
	mu     sync.Mutex
	amount storage.AccountBalance
}

// New returns a StorageInMemory.
func New() *StorageInMemory {
	return &StorageInMemory{
		m: sync.Map{},
	}
}

// Set sets a new account with balance.
func (s *StorageInMemory) Set(account storage.AccountID, amount storage.AccountBalance) error {
	if account == "" {
		return ErrAccNameInvalid
	}
	body := &balance{
		mu:     sync.Mutex{},
		amount: amount,
	}
	s.m.Store(account, body)
	return nil
}

// Get reads a balance of account.
func (s *StorageInMemory) Get(account storage.AccountID) (storage.AccountBalance, error) {
	balance, err := s.getBalance(account)
	if err != nil {
		return 0, err
	}
	return balance.amount, nil
}

func (s *StorageInMemory) getBalance(account storage.AccountID) (*balance, error) {
	body, ok := s.m.Load(account)
	if !ok {
		return nil, ErrAccIsAbsent
	}
	return body.(*balance), nil
}

// Transfer ransfers money from accountFrom to accountTo using the amount.
func (s *StorageInMemory) Transfer(accountFrom, accountTo storage.AccountID, amountToTransfer storage.AccountBalance) error {
	balanceFrom, err := s.getBalance(accountFrom)
	if err != nil {
		return err
	}
	balanceTo, err := s.getBalance(accountTo)
	if err != nil {
		return err
	}

	balanceFrom.mu.Lock()
	balanceTo.mu.Lock()
	defer balanceFrom.mu.Unlock()
	defer balanceTo.mu.Unlock()

	if balanceFrom.amount < amountToTransfer {
		return ErrNotEnoughBalance
	}
	// todo del (для отладки)
	// fmt.Println("операция: ", balanceFrom.amount, balanceTo.amount, balanceFrom.amount+balanceTo.amount)
	balanceFrom.amount -= amountToTransfer
	balanceTo.amount += amountToTransfer
	return nil
}

// Print метод только для отладки todo удалить реализацию и из интерфейса storage.Storage.
func (s *StorageInMemory) Print() {
	mapToPrint := map[string]interface{}{}
	s.m.Range(func(key, value interface{}) bool {
		mapToPrint[fmt.Sprint(key)] = value.(*balance).amount
		return true
	})
	fmt.Println("Текущее состояние счетов: ", mapToPrint)
}
