package internal

type Storage interface {
	UserByLogin(login string) (*User, error)
	AddUser(user User) error
	AddOrder(order Order) error
	OrdersByUser(login string) ([]Order, error)
	UpdateUserBalance(login string, accrual float64) error
	AddTransaction(transaction Transaction) (*Transaction, error)
	TransactionsByUser(login string) ([]Transaction, error)
}
