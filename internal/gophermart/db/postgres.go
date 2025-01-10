package db

import (
	"database/sql"
	"embed"
	"log/slog"
	"musthave_tpl/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *slog.Logger
	pSQL   sq.StatementBuilderType
}
type GeneralRepository interface {
	OrdersByUser(string) ([]models.Order, error)
	AddOrder(order models.Order) error
	AddUser(user models.User) error
	UserByLogin(login string) (*models.User, error)
	AddTransaction(transaction models.Transaction) (rows int64, err error)
	TransactionsByUser(login string) ([]models.Transaction, error)
}

type UserRepository interface {
	AddUser(user models.User) error
	UserByLogin(login string) (*models.User, error)
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewPostgresStorage(db *sql.DB, logger slog.Handler) (*PostgresRepository, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
	ps := &PostgresRepository{
		logger: slog.New(logger),
		db:     db,
		pSQL:   psql,
	}
	err := ps.init()
	return ps, err
}

func (ps *PostgresRepository) init() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(ps.db, "migrations"); err != nil {
		return err
	}
	return nil
}
