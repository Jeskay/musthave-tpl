package db

import (
	"database/sql"
	"log/slog"
	"musthave_tpl/internal"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresStorage(db *sql.DB, logger slog.Handler) (*PostgresStorage, error) {
	ps := &PostgresStorage{
		logger: slog.New(logger),
		db:     db,
	}
	return ps, nil
}

func (ps *PostgresStorage) init() error {
	query := `
		CREATE TABLE IF NOT EXIST user (
			login varchar(500) PRIMARY KEY UNIQUE,
			password text,
			balance bigint
		);
	`
	_, err := ps.db.Exec(query)
	return err
}

func (ps *PostgresStorage) AddUser(user internal.User) error {
	query := `
		INSERT INTO user (name, password)
		VALUES ($1, $2);
	`
	_, err := ps.db.Exec(query, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) UserByLogin(login string) (*internal.User, error) {
	var user internal.User
	query := `
		SELECT 
			* 
		FROM user 
		WHERE login = $1
	`
	row := ps.db.QueryRow(query, login)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.Scan(&user.Login, &user.Password, &user.Balance); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ps *PostgresStorage) OrdersByUser(login string) ([]internal.Order, error) {
	return nil, nil
}

func (ps *PostgresStorage) UpdateUser(user internal.User) (*internal.User, error) {
	return nil, nil
}

func (ps *PostgresStorage) AddTransaction(user internal.User, transaction internal.Transaction) (*internal.Transaction, error) {
	return nil, nil
}

func (ps *PostgresStorage) TransactionsByUser(user internal.User) ([]internal.Transaction, error) {
	return nil, nil
}
