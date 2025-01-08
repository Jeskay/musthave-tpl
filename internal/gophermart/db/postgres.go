package db

import (
	"database/sql"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *slog.Logger
	pSQL   sq.StatementBuilderType
}

func NewPostgresStorage(db *sql.DB, logger slog.Handler) (*PostgresStorage, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
	ps := &PostgresStorage{
		logger: slog.New(logger),
		db:     db,
		pSQL:   psql,
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
			uploaded_at timestamptz DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS withdrawals (
			id SERIAL PRIMARY KEY,
			user_login varchar(500) REFERENCES users (login),
			order_id bigint,
			amount double precision,
			processed_at timestamptz DEFAULT NOW()
		);
	`
	_, err := ps.db.Exec(query)
	return err
}
