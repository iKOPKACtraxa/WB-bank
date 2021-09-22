package storage

type (
	AccountID      string
	AccountBalance int
)

type Storage interface {
	Set(account AccountID, balance AccountBalance) error
	Get(account AccountID) (AccountBalance, error)
	Transfer(accountFrom, accountTo AccountID, Balance AccountBalance) error
	Print() // todo только для отладки, удалить
}
