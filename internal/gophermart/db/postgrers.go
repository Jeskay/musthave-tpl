package db

import (
	"context"
	"database/sql"
	"log/slog"
	"musthave_tpl/internal"
	"musthave_tpl/internal/gophermart"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	err := ps.init()
	return ps, err
}

func (ps *PostgresStorage) init() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			login varchar(500) PRIMARY KEY UNIQUE,
			password text,
			balance double precision DEFAULT 0,
			withdrawn double precision DEFAULT 0
		);
		CREATE TABLE IF NOT EXISTS orders (
			user_login varchar(500) REFERENCES users (login),
			id bigint PRIMARY KEY UNIQUE,
			status varchar(200),
			accrual double precision,
			uploaded_at timestamptz
		);
		CREATE TABLE IF NOT EXISTS withdrawals (
			id SERIAL PRIMARY KEY,
			user_login varchar(500) REFERENCES users (login),
			order_id bigint,
			amount double precision,
			processed_at timestamptz
		);
	`
	_, err := ps.db.Exec(query)
	return err
}

func (ps *PostgresStorage) AddUser(user internal.User) error {
	row := ps.db.QueryRow(`SELECT COUNT(*) FROM users WHERE login = $1`, user.Login)
	if row.Err() != nil {
		return row.Err()
	}
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count != 0 {
		return &gophermart.UsedLoginError{}
	}
	query := `
		INSERT INTO users (login, password)
		VALUES ($1, $2);
	`
	_, err = ps.db.Exec(query, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) UserByLogin(login string) (*internal.User, error) {
	var (
		userLogin     string
		userPassword  string
		userBalance   sql.NullFloat64
		userWithdrawn sql.NullFloat64
	)
	query := `
		SELECT 
			login,
			password, 
			balance,
			withdrawn
		FROM users 
		WHERE login = $1
	`
	row := ps.db.QueryRow(query, login)
	if row.Err() != nil {
		return nil, row.Err()
	}
	if err := row.Scan(&userLogin, &userPassword, &userBalance, &userWithdrawn); err != nil {
		return nil, err
	}
	return &internal.User{
		Login:     userLogin,
		Password:  userPassword,
		Balance:   userBalance.Float64,
		Withdrawn: userWithdrawn.Float64,
	}, nil
}

func (ps *PostgresStorage) OrdersByUser(login string) ([]internal.Order, error) {
	query := `
		SELECT
			id,
			user_login,
			status,
			accrual,
			uploaded_at
		FROM orders
		WHERE user_login = $1
		ORDER BY uploaded_at ASC;
	`
	rows, err := ps.db.Query(query, login)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	var orders []internal.Order
	for rows.Next() {
		var (
			id         int64
			userLogin  string
			status     string
			accrual    float64
			uploadedAt time.Time
		)
		err := rows.Scan(&id, &userLogin, &status, &accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, internal.Order{
			User:       internal.User{Login: userLogin},
			Number:     id,
			Status:     internal.OrderStatus(status),
			Accrual:    accrual,
			UploadedAt: uploadedAt,
		})
	}
	return orders, nil
}

func (ps *PostgresStorage) AddOrder(order internal.Order) error {
	rows, err := ps.db.Query(`SELECT (user_login) FROM orders WHERE id = $1`, order.Number)
	if err != nil {
		return err
	}
	if rows.Err() != nil {
		return err
	}
	if rows.Next() {
		var login string
		err := rows.Scan(&login)
		if err != nil {
			return err
		}
		if login == order.User.Login {
			return &gophermart.OrderExists{}
		}
		return &gophermart.OrderUsed{}
	}
	ctx := context.Background()
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO orders (id, user_login, status, accrual, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5);
	`
	_, err = ps.db.ExecContext(ctx, query, order.Number, order.User.Login, order.Status, order.Accrual, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}
	query2 := `
		UPDATE users
		SET balance = balance + $1
		WHERE login = $2;
	`
	_, err = ps.db.ExecContext(ctx, query2, order.Accrual, order.User.Login)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (ps *PostgresStorage) UpdateUserBalance(login string, accrual float64) error {
	query := `
		UPDATE users
		SET balance = balance + $1
		WHERE login = $2;
	`
	row := ps.db.QueryRow(query, accrual, login)
	return row.Err()
}

func (ps *PostgresStorage) AddTransaction(transaction internal.Transaction) (*internal.Transaction, error) {
	ctx := context.Background()
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	query1 := `
	INSERT INTO withdrawals (user_login, order_id, amount, processed_at)
	VALUES ($1, $2, $3, $4);
	`
	_, err = tx.ExecContext(ctx, query1, transaction.User, transaction.ID, transaction.Amount, transaction.Date)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	query2 := `
	UPDATE users
	SET balance = balance - $1,
	withdrawn = withdrawn + $1
	WHERE login = $2
	`
	_, err = tx.ExecContext(ctx, query2, transaction.Amount, transaction.User)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	return nil, err
}

func (ps *PostgresStorage) TransactionsByUser(login string) ([]internal.Transaction, error) {
	query := `
	SELECT 
		user_login,
		order_id,
		amount,
		processed_at
	 FROM withdrawals
	WHERE user_login = $1
	`
	rows, err := ps.db.Query(query, login)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	var transactions []internal.Transaction
	for rows.Next() {
		var (
			login  string
			id     int64
			amount float64
			date   time.Time
		)
		err := rows.Scan(&login, &id, &amount, &date)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, internal.Transaction{ID: id, Amount: amount, User: login, Date: date})
	}
	return transactions, nil
}
