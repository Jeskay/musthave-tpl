package db

import (
	"database/sql"
	"embed"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *slog.Logger
	pSQL   sq.StatementBuilderType
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

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
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(ps.db, "migrations"); err != nil {
		return err
	}
	return nil
}
