package internal

type Storage interface {
	UserByLogin(login string) (*User, error)
	AddUser(user User) error
	OrdersByUser(login string) ([]Order, error)
	UpdateUser(user User) (*User, error)
	AddTransaction(user User, transaction Transaction) (*Transaction, error)
	TransactionsByUser(user User) ([]Transaction, error)
}
