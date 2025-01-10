package db

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"
	"musthave_tpl/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db        *sql.DB
	logger    *slog.Logger
	pSQL      sq.StatementBuilderType
	txTimeout time.Duration
}
type GeneralRepository interface {
	OrdersByUser(ctx context.Context, login string) ([]models.Order, error)
	AddOrder(ctx context.Context, order models.Order) error
	AddUser(ctx context.Context, user models.User) error
	UserByLogin(ctx context.Context, login string) (*models.User, error)
	AddTransaction(ctx context.Context, transaction models.Transaction) (rows int64, err error)
	TransactionsByUser(ctx context.Context, login string) ([]models.Transaction, error)
}

type UserRepository interface {
	AddUser(ctx context.Context, user models.User) error
	UserByLogin(ctx context.Context, login string) (*models.User, error)
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewPostgresStorage(db *sql.DB, logger slog.Handler) (*PostgresRepository, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(db)
	ps := &PostgresRepository{
		logger:    slog.New(logger),
		db:        db,
		pSQL:      psql,
		txTimeout: time.Second * 5,
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
