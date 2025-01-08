package internal

import "musthave_tpl/internal/models"

type Storage interface {
	UserByLogin(login string) (*models.User, error)
	AddUser(user models.User) error
	AddOrder(order models.Order) error
	OrdersByUser(login string) ([]models.Order, error)
	AddTransaction(transaction models.Transaction) (int64, error)
	TransactionsByUser(login string) ([]models.Transaction, error)
}
