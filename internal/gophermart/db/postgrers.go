package db

import (
	"context"
	"database/sql"
	"log/slog"
	"musthave_tpl/internal"
	"musthave_tpl/internal/gophermart"
	"time"

	sq "github.com/Masterminds/squirrel"

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

func (ps *PostgresStorage) UserExist(user internal.User) (bool, error) {
	row := ps.db.QueryRow(`SELECT COUNT(*) FROM users WHERE login = $1`, user.Login)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
	return true, nil
}

func (ps *PostgresStorage) AddUser(user internal.User) error {
	exist, err := ps.UserExist(user)
	if err != nil {
		return err
	} else if !exist {
		return &gophermart.UsedLoginError{}
	}
	queryAddUser := sq.Insert(
		"users",
	).Columns(
		"login",
		"password",
	).Values(user.Login, user.Password)

	_, err = queryAddUser.RunWith(ps.db).Exec()
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
	query := sq.Select(
		"login", "password", "balance", "withdrawn",
	).From(
		"users",
	).Where(sq.Eq{"login": login})

	row := query.RunWith(ps.db).QueryRow()
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
	query := sq.Select(
		"id", "user_login", "status", "accrual", "uploaded_at",
	).From(
		"orders",
	).Where(
		sq.Eq{"user_login": login},
	).OrderBy("uploaded_at ASC")

	rows, err := query.RunWith(ps.db).Query()
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
	queryUser := sq.Select(
		"user_login",
	).From(
		"orders",
	).Where(sq.Eq{"id": order.Number})

	rows, err := queryUser.RunWith(ps.db).Query()
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

	orderQuery := sq.Insert(
		"orders",
	).Columns(
		"id",
		"user_login",
		"status",
		"accrual",
		"uploaded_at",
	).Values(order.Number, order.User.Login, order.Status, order.Accrual, time.Now())

	_, err = orderQuery.RunWith(ps.db).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	balanceQuery := sq.Update(
		"update",
	).Set(
		"balance",
		sq.Expr("balance + ?", order.Accrual),
	).Where(
		sq.Eq{"login": order.User.Login},
	)

	_, err = balanceQuery.RunWith(ps.db).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (ps *PostgresStorage) AddTransaction(transaction internal.Transaction) (*internal.Transaction, error) {
	ctx := context.Background()
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryWithdraw := sq.Insert(
		"withdrawals",
	).Columns(
		"user_login",
		"order_id",
		"amount",
		"processed_at",
	).Values(transaction.User, transaction.ID, transaction.Amount, transaction.Date)

	_, err = queryWithdraw.RunWith(ps.db).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	queryUser := sq.Update(
		"users",
	).Set(
		"balance",
		sq.Expr("balance - ?", transaction.Amount),
	).Set(
		"withdrawn",
		sq.Expr("withdrawn + ?", transaction.Amount),
	).Where(sq.Eq{"login": transaction.User})
	_, err = queryUser.RunWith(ps.db).ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	return nil, err
}

func (ps *PostgresStorage) TransactionsByUser(login string) ([]internal.Transaction, error) {
	query := sq.Select(
		"user_login",
		"order_id",
		"amount",
		"processed_at",
	).From(
		"withdrawals",
	).Where(sq.Eq{"user_login": login})

	rows, err := query.RunWith(ps.db).Query()
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
